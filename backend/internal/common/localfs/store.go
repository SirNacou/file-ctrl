package localfs

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type Store struct {
	baseDir string
	mu      sync.Mutex
}

func NewStore(baseDir string) *Store { return &Store{baseDir: baseDir} }

func (s *Store) StorageDir() string { return s.baseDir }

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
) ([]domain.FileMeta, error) {
	root, err := s.openRoot(sourcePath)
	if err != nil {
		return nil, err
	}
	defer closeLog(ctx, root, sourcePath.String())

	d, err := root.Open(dir.FS())
	if err != nil {
		return nil, err
	}
	defer closeLog(ctx, d, dir.Display())

	entries, err := d.ReadDir(-1) // fails if not a directory
	if err != nil {
		return nil, err
	}

	out := make([]domain.FileMeta, 0, len(entries))
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue // removed between ReadDir and Info
		}
		out = append(out, domain.FileMeta{
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
	s.mu.Lock()
	defer s.mu.Unlock()
	root, err := s.openRoot(sourcePath)
	if err != nil {
		return err
	}
	defer closeLog(ctx, root, sourcePath.String())

	if err := root.Mkdir(target.FS(), 0o755); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return err
		}
		return err
	}

	return nil
}

func (s *Store) Delete(
	ctx context.Context,
	sourcePath domain.SourcePath,
	target domain.EntryPath,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	root, err := s.openRoot(sourcePath)
	if err != nil {
		return err
	}
	defer closeLog(ctx, root, target.Display())

	if err := root.Remove(target.FS()); err != nil {
		return err
	}
	return nil
}

func (s *Store) Rename(
	ctx context.Context,
	sourcePath domain.SourcePath,
	from, to domain.EntryPath,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	root, err := s.openRoot(sourcePath)
	if err != nil {
		return domain.ErrStorageSourceNotFound
	}
	defer closeLog(ctx, root, sourcePath.String())

	if _, err := root.Lstat(from.FS()); err != nil {
		return err
	}
	// Rename silently overwrites on Unix, so refuse if the target exists.
	if _, err := root.Lstat(to.FS()); err == nil {
		return err
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	if err := root.Rename(from.FS(), to.FS()); err != nil {
		return err
	}
	return nil
}

func (s *Store) OpenFile(
	ctx context.Context,
	sourcePath domain.SourcePath,
	target domain.EntryPath,
) (*os.File, domain.FileMeta, error) {
	root, err := s.openRoot(sourcePath)
	if err != nil {
		return nil, domain.FileMeta{}, err
	}
	defer closeLog(ctx, root, sourcePath.String())

	file, err := root.Open(target.FS())
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, domain.FileMeta{}, err
		}
		return nil, domain.FileMeta{}, err
	}

	info, err := file.Stat()
	if err != nil {
		closeLog(ctx, file, target.FS())
		return nil, domain.FileMeta{}, err
	}

	if info.IsDir() {
		closeLog(ctx, file, target.FS())
		return nil, domain.FileMeta{}, err
	}

	return file, domain.FileMeta{
		Name:    info.Name(),
		ModTime: info.ModTime(),
		Size:    info.Size(),
		IsDir:   info.IsDir(),
	}, nil
}

func (s *Store) OpenFS(sourcePath domain.SourcePath) (fs.FS, io.Closer, error) {
	root, err := s.openRoot(sourcePath)
	if err != nil {
		return nil, nil, err
	}
	return root.FS(), root, nil // caller closes: root must outlive the zip stream
}

func (s *Store) CommitStagedFile(
	stagingFilePath string,
	sourcePath domain.SourcePath,
	targetPath domain.EntryPath,
) error {
	root, err := s.openRoot(sourcePath)
	if err != nil {
		return err
	}
	defer root.Close()

	parentDir := filepath.Dir(targetPath.FS())
	if parentDir != "." && parentDir != "/" {
		_ = root.MkdirAll(parentDir, 0o755)
	}

	srcFile, err := os.Open(stagingFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := root.OpenFile(targetPath.FS(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}

func (s *Store) openRoot(relPath domain.SourcePath) (*os.Root, error) {
	base, err := os.OpenRoot(s.baseDir)
	if err != nil {
		return nil, err
	}
	defer base.Close()

	return base.OpenRoot(filepath.FromSlash(relPath.String()))
}

func closeLog(ctx context.Context, c io.Closer, target string) {
	if err := c.Close(); err != nil {
		slog.WarnContext(ctx, "failed to close file",
			slog.String("target", target),
			slog.Any("err", err),
		)
	}
}
