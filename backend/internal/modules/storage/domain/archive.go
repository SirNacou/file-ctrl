package domain

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
)

type archiveEntry struct {
	src   string // path inside fsys, used to open the file
	dst   string // path inside the zip
	isDir bool
}

type Archive struct {
	fsys    fs.FS
	closer  io.Closer
	name    ArchiveName
	entries []archiveEntry
}

func NewArchive(
	fsys fs.FS,
	closer io.Closer,
	name ArchiveName,
	paths []EntryPath,
) (*Archive, error) {
	if len(paths) == 0 {
		return nil, ErrEmptySelection
	}

	a := &Archive{fsys: fsys, closer: closer, name: name}
	seen := map[string]bool{} // zip paths already added
	single := len(paths) == 1

	for _, p := range paths {
		if err := a.collect(p.FS(), single, seen); err != nil {
			return nil, toDomainError(err)
		}
	}
	if len(a.entries) == 0 {
		return nil, ErrEmptySelection
	}
	return a, nil
}

// collect walks one selected path and records every file/dir under it.
func (a *Archive) collect(root string, single bool, seen map[string]bool) error {
	prefix, err := a.zipPrefix(root, single)
	if err != nil {
		return err
	}

	return fs.WalkDir(a.fsys, root, func(src string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && !d.Type().IsRegular() { // skip symlinks, devices, etc.
			return nil
		}

		dst := path.Join(prefix, relativeTo(root, src))
		if dst == "." || seen[dst] { // "." = the folder itself (flattened)
			return nil
		}
		seen[dst] = true
		a.entries = append(a.entries, archiveEntry{src: src, dst: dst, isDir: d.IsDir()})
		return nil
	})
}

// zipPrefix returns the folder name that everything from this selection
// is placed under inside the zip ("" means directly at the zip root).
func (a *Archive) zipPrefix(root string, single bool) (string, error) {
	if root == "." { // whole storage selected
		return "", nil
	}
	if single {
		info, err := fs.Stat(a.fsys, root)
		if err != nil {
			return "", err
		}
		if info.IsDir() { // one folder: use its contents as the zip root
			return "", nil
		}
	}
	return path.Base(root), nil
}

// relativeTo returns src relative to root ("." if src is root itself).
func relativeTo(root, src string) string {
	switch {
	case root == ".":
		return src
	case src == root:
		return "."
	default:
		return strings.TrimPrefix(src, root+"/")
	}
}

func toDomainError(err error) error {
	if errors.Is(err, fs.ErrNotExist) {
		return ErrNotFound
	}
	return ErrInvalidPath
}

func (a *Archive) Name() string { return a.name.String() }

func (a *Archive) Close() error { return a.closer.Close() }

func (a *Archive) WriteZip(ctx context.Context, w io.Writer) error {
	zw := zip.NewWriter(w)
	for _, e := range a.entries {
		if err := ctx.Err(); err != nil {
			return err
		}
		var err error
		if e.isDir {
			err = addDir(zw, e)
		} else {
			err = a.addFile(zw, e)
		}
		if err != nil {
			return fmt.Errorf("add %q: %w", e.src, err)
		}
	}
	return zw.Close()
}

func addDir(zw *zip.Writer, e archiveEntry) error {
	_, err := zw.CreateHeader(&zip.FileHeader{Name: e.dst + "/", Method: zip.Store})
	return err
}

func (a *Archive) addFile(zw *zip.Writer, e archiveEntry) error {
	f, err := a.fsys.Open(e.src)
	if err != nil {
		return err
	}
	defer f.Close()

	w, err := zw.Create(e.dst)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, f)
	return err
}
