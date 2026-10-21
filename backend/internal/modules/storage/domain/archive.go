package domain

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
)

type Archive struct {
	fsys   fs.FS
	closer io.Closer
	name   ArchiveName
	files  []string
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
	seen := map[string]struct{}{}
	a := &Archive{fsys: fsys, closer: closer, name: name}

	for _, p := range paths {
		err := fs.WalkDir(fsys, p.FS(), func(name string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.Type().IsRegular() { // dirs are walked, symlinks/devices skipped
				return nil
			}
			if _, dup := seen[name]; !dup {
				seen[name] = struct{}{}
				a.files = append(a.files, name)
			}
			return nil
		})
		switch {
		case errors.Is(err, fs.ErrNotExist):
			return nil, ErrNotFound
		case err != nil:
			return nil, ErrInvalidPath
		}
	}
	if len(a.files) == 0 {
		return nil, ErrEmptySelection
	}
	return a, nil
}

func (a *Archive) Name() string {
	return a.name.String()
}

func (a *Archive) Close() error { return a.closer.Close() }

func (a *Archive) WriteZip(ctx context.Context, w io.Writer) error {
	zw := zip.NewWriter(w)
	for _, file := range a.files {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := a.add(zw, file); err != nil {
			return fmt.Errorf("add %q: %w", file, err)
		}
	}
	return zw.Close()
}

func (a *Archive) add(zw *zip.Writer, name string) error {
	f, err := a.fsys.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}
	h, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	h.Name, h.Method = name, zip.Deflate

	dst, err := zw.CreateHeader(h)
	if err != nil {
		return err
	}
	_, err = io.Copy(dst, f)
	return err
}
