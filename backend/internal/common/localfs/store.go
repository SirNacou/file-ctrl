package localfs

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type Store struct {
	baseDir string
}

func NewStore(baseDir string) *Store { return &Store{baseDir: baseDir} }

func (s *Store) DirExists(ctx context.Context, relPath string) (bool, error) {
	base, err := os.OpenRoot(s.baseDir)
	if err != nil {
		return false, err
	}
	defer closeLog(ctx, base, s.baseDir)

	info, err := base.Stat(filepath.FromSlash(relPath))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return info.IsDir(), nil
}

func (s *Store) ListDirs() ([]string, error) {
	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, err
	}

	dirs := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		}
	}
	return dirs, nil
}

func (s *Store) ReadDir(
	ctx context.Context,
	sourcePath domain.SourcePath,
	dir domain.EntryPath,
) ([]domain.FileEntry, error) {
	root, err := s.openRoot(sourcePath.String())
	if err != nil {
		return nil, domain.ErrStorageSourceNotFound
	}
	defer closeLog(ctx, root, sourcePath.String())

	d, err := root.Open(dir.FS())
	if err != nil {
		return nil, mapFSError(err)
	}
	defer closeLog(ctx, d, dir.Display())

	entries, err := d.ReadDir(-1) // fails if not a directory
	if err != nil {
		return nil, domain.ErrInvalidPath
	}

	out := make([]domain.FileEntry, 0, len(entries))
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue // removed between ReadDir and Info
		}
		out = append(out, domain.FileEntry{
			Name:    e.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			IsDir:   e.IsDir(),
		})
	}
	return out, nil
}

func (s *Store) Mkdir(
	ctx context.Context,
	sourcePath domain.SourcePath,
	target domain.EntryPath,
) error {
	root, err := s.openRoot(sourcePath.String())
	if err != nil {
		return domain.ErrStorageSourceNotFound
	}
	defer closeLog(ctx, root, sourcePath.String())

	if err := root.Mkdir(target.FS(), 0o755); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return domain.ErrAlreadyExists
		}
		return domain.ErrInvalidPath
	}

	return nil
}

func (s *Store) Delete(
	ctx context.Context,
	sourcePath domain.SourcePath,
	target domain.EntryPath,
) error {
	root, err := s.openRoot(sourcePath.String())
	if err != nil {
		return domain.ErrStorageSourceNotFound
	}
	defer closeLog(ctx, root, target.Display())

	if err := root.Remove(target.FS()); err != nil {
		return mapFSError(err)
	}
	return nil
}

func (s *Store) Rename(
	ctx context.Context,
	sourcePath domain.SourcePath,
	from, to domain.EntryPath,
) error {
	root, err := s.openRoot(sourcePath.String())
	if err != nil {
		return domain.ErrStorageSourceNotFound
	}
	defer closeLog(ctx, root, sourcePath.String())

	if _, err := root.Lstat(from.FS()); err != nil {
		return mapFSError(err)
	}
	// Rename silently overwrites on Unix, so refuse if the target exists.
	if _, err := root.Lstat(to.FS()); err == nil {
		return domain.ErrAlreadyExists
	} else if !errors.Is(err, fs.ErrNotExist) {
		return domain.ErrInvalidPath
	}
	if err := root.Rename(from.FS(), to.FS()); err != nil {
		return mapFSError(err)
	}
	return nil
}

func mapFSError(err error) error {
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return domain.ErrNotFound
	case errors.Is(err, fs.ErrExist):
		return domain.ErrAlreadyExists
	default:
		return domain.ErrInvalidPath
	}
}

func (s *Store) openRoot(relPath string) (*os.Root, error) {
	base, err := os.OpenRoot(s.baseDir)
	if err != nil {
		return nil, err
	}
	defer base.Close()

	return base.OpenRoot(filepath.FromSlash(relPath))
}

func closeLog(ctx context.Context, c io.Closer, target string) {
	if err := c.Close(); err != nil {
		slog.WarnContext(ctx, "failed to close os.Root",
			slog.String("target", target),
			slog.Any("err", err),
		)
	}
}
