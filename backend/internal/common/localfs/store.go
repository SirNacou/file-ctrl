package localfs

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type Store struct {
	baseDir string
}

func NewStore(baseDir string) *Store { return &Store{baseDir: baseDir} }

func (s *Store) OpenSource(sourcePath string) (*os.Root, error) {
	base, err := os.OpenRoot(s.baseDir)
	if err != nil {
		return nil, err
	}
	defer base.Close()

	return base.OpenRoot(filepath.FromSlash(sourcePath))
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

func CloseLog(ctx context.Context, c io.Closer, what string) {
	if err := c.Close(); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("failed to close %s", what), slog.Any("error", err))
	}
}
