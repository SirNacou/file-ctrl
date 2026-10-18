package app

import (
	"context"
	"os"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type SourceOpener interface {
	OpenSource(sourcePath string) (*os.Root, error)
}
type StorageDirLister interface {
	ListDirs() ([]string, error)
}

type StorageSourcesInserter interface {
	Insert(ctx context.Context, model *domain.StorageSource) error
}

type StorageSourcesFinder interface {
	Find(ctx context.Context, path string) (*domain.StorageSource, error)
}
type StorageSourcesUpdater interface {
	Update(ctx context.Context, model *domain.StorageSource) error
}
