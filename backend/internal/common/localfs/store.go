package localfs

import (
	"os"
	"path/filepath"
)

type SourceStore struct {
	baseDir string
}

func NewSourceStore(baseDir string) *SourceStore { return &SourceStore{baseDir: baseDir} }

func (s *SourceStore) OpenSource(sourcePath string) (*os.Root, error) {
	base, err := os.OpenRoot(s.baseDir)
	if err != nil {
		return nil, err
	}
	defer base.Close()

	return base.OpenRoot(filepath.FromSlash(sourcePath))
}

func (s *SourceStore) ListDirs() ([]string, error) {
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
