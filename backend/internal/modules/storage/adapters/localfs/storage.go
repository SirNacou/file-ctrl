package localfs

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type Storage struct {
	baseDir string
	mu      sync.Mutex
}

func NewStorage(baseDir string) *Storage { return &Storage{baseDir: baseDir} }

func (s *Storage) StorageDir() string { return s.baseDir }

func (s *Storage) DirExists(ctx context.Context, relPath string) (bool, error) {
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

func (s *Storage) ListDirs() ([]string, error) {
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

func (s *Storage) ReadDir(
	ctx context.Context,
	sourcePath domain.SourcePath,
	dir domain.EntryPath,
) ([]domain.FileMeta, error) {
	root, err := s.openRoot(ctx, sourcePath)
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
		out = append(out, toFileMeta(info))
	}
	return out, nil
}

func (s *Storage) Mkdir(
	ctx context.Context,
	sourcePath domain.SourcePath,
	target domain.EntryPath,
) error {
	return s.withWriteRoot(ctx, sourcePath, func(root *os.Root) error {
		return root.Mkdir(target.FS(), 0o755)
	})
}

func (s *Storage) Delete(
	ctx context.Context,
	sourcePath domain.SourcePath,
	target domain.EntryPath,
) error {
	return s.withWriteRoot(ctx, sourcePath, func(root *os.Root) error {
		return root.RemoveAll(target.FS())
	})
}

func (s *Storage) Rename(
	ctx context.Context,
	sourcePath domain.SourcePath,
	from, to domain.EntryPath,
) error {
	return s.withWriteRoot(ctx, sourcePath, func(root *os.Root) error {
		if _, err := root.Lstat(from.FS()); err != nil {
			return err
		}
		// Rename silently overwrites on Unix, so refuse if the target exists.
		if _, err := root.Lstat(to.FS()); err == nil {
			return fs.ErrExist
		} else if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		return root.Rename(from.FS(), to.FS())
	})
}

func (s *Storage) OpenFile(
	ctx context.Context,
	sourcePath domain.SourcePath,
	target domain.EntryPath,
) (*os.File, domain.FileMeta, error) {
	root, err := s.openRoot(ctx, sourcePath)
	if err != nil {
		return nil, domain.FileMeta{}, err
	}
	defer closeLog(ctx, root, sourcePath.String())

	file, err := root.Open(target.FS())
	if err != nil {
		return nil, domain.FileMeta{}, err
	}

	info, err := file.Stat()
	if err != nil {
		closeLog(ctx, file, target.FS())
		return nil, domain.FileMeta{}, err
	}

	if info.IsDir() {
		closeLog(ctx, file, target.FS())
		return nil, domain.FileMeta{}, fs.ErrInvalid
	}

	return file, toFileMeta(info), nil
}

func (s *Storage) OpenFS(sourcePath domain.SourcePath) (fs.FS, io.Closer, error) {
	root, err := s.openRoot(context.Background(), sourcePath)
	if err != nil {
		return nil, nil, err
	}
	return root.FS(), root, nil // caller closes: root must outlive the zip stream
}

func (s *Storage) CommitStagedFile(
	ctx context.Context,
	stagingFilePath string,
	sourcePath domain.SourcePath,
	targetPath domain.EntryPath,
) error {
	// 1. Open staging file first (before touching destination root)
	srcFile, err := os.Open(stagingFilePath)
	if err != nil {
		return err
	}
	defer closeLog(ctx, srcFile, stagingFilePath)

	root, err := s.openRoot(ctx, sourcePath)
	if err != nil {
		return err
	}
	defer closeLog(ctx, root, sourcePath.String())

	// 2. Stream bytes to a temporary file at the root of the storage source
	// BEFORE creating any target folders or locking the mutex.
	// If io.Copy fails, no redundant folders exist on disk!
	tmpName := ".upload-" + filepath.Base(stagingFilePath) + ".tmp"
	dstFile, err := root.OpenFile(tmpName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	committed := false
	defer func() {
		if !committed {
			_ = root.Remove(tmpName)
		}
	}()

	_, copyErr := io.Copy(dstFile, srcFile)
	closeErr := dstFile.Close()
	if err := errors.Join(copyErr, closeErr); err != nil {
		return err
	}

	// 3. Now that all file bytes are safely on disk, lock the mutex,
	// auto-create the nested folder path, and atomically move the file in place.
	s.mu.Lock()
	defer s.mu.Unlock()

	targetFS := targetPath.FS()
	parentDir := filepath.Dir(targetFS)

	undoDirs, err := mkdirAllWithRollback(root, parentDir)
	if err != nil {
		return err
	}

	if err := root.Rename(tmpName, targetFS); err != nil {
		undoDirs() // Undo any newly created folders if Rename fails
		return err
	}

	committed = true
	return nil
}

// mkdirAllWithRollback creates only the missing directories in relDir and returns
// an undo function that removes those exact newly created directories in reverse order.
func mkdirAllWithRollback(root *os.Root, relDir string) (func(), error) {
	if relDir == "." || relDir == "/" || relDir == "" {
		return func() {}, nil
	}

	var created []string
	parts := strings.Split(filepath.ToSlash(relDir), "/")
	curr := ""

	for _, part := range parts {
		if part == "" || part == "." {
			continue
		}
		if curr == "" {
			curr = part
		} else {
			curr = curr + "/" + part
		}

		fsPath := filepath.FromSlash(curr)
		info, statErr := root.Stat(fsPath)
		if statErr == nil {
			if !info.IsDir() {
				// A regular file already exists where a folder needs to be!
				for _, c := range slices.Backward(created) {
					_ = root.Remove(c)
				}
				return nil, fs.ErrExist
			}
			continue // Folder already existed; don't add to rollback list
		}
		if !errors.Is(statErr, fs.ErrNotExist) {
			for _, c := range slices.Backward(created) {
				_ = root.Remove(c)
			}
			return nil, statErr
		}

		if mkdirErr := root.Mkdir(fsPath, 0o755); mkdirErr != nil {
			for _, c := range slices.Backward(created) {
				_ = root.Remove(c)
			}
			return nil, mkdirErr
		}
		created = append(created, fsPath)
	}

	return func() {
		for _, c := range slices.Backward(created) {
			_ = root.Remove(c)
		}
	}, nil
}

func (s *Storage) withWriteRoot(
	ctx context.Context,
	sourcePath domain.SourcePath,
	fn func(root *os.Root) error,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	root, err := s.openRoot(ctx, sourcePath)
	if err != nil {
		return err
	}
	defer closeLog(ctx, root, sourcePath.String())

	return fn(root)
}

func (s *Storage) openRoot(ctx context.Context, relPath domain.SourcePath) (*os.Root, error) {
	base, err := os.OpenRoot(s.baseDir)
	if err != nil {
		return nil, err
	}
	defer closeLog(ctx, base, s.baseDir)

	return base.OpenRoot(filepath.FromSlash(relPath.String()))
}

func toFileMeta(info fs.FileInfo) domain.FileMeta {
	return domain.FileMeta{
		Name:    info.Name(),
		Size:    info.Size(),
		ModTime: info.ModTime(),
		IsDir:   info.IsDir(),
	}
}

func closeLog(ctx context.Context, c io.Closer, target string) {
	if err := c.Close(); err != nil {
		slog.WarnContext(ctx, "failed to close file",
			slog.String("target", target),
			slog.Any("err", err),
		)
	}
}
