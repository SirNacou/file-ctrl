package app

import (
	"context"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type StorageSourcesInserter interface {
	Insert(ctx context.Context, model *domain.StorageSource) error
}

type StorageSourcesFinder interface {
	Find(ctx context.Context, path string) (*domain.StorageSource, error)
}
type StorageSourcesUpdater interface {
	Update(ctx context.Context, model *domain.StorageSource) error
}

type FSDirReader interface {
	ReadDir(
		ctx context.Context,
		sourcePath domain.SourcePath,
		dir domain.EntryPath,
	) ([]domain.FileEntry, error)
}

type FSDirExistsChecker interface {
	DirExists(ctx context.Context, relPath string) (bool, error)
}

type FSDirCreator interface {
	Mkdir(ctx context.Context, sourcePath domain.SourcePath, target domain.EntryPath) error
}

type FSFileRenamer interface {
	Rename(ctx context.Context, sourcePath domain.SourcePath, from, to domain.EntryPath) error
}

type FSDirLister interface {
	ListDirs() ([]string, error)
}

type FSFileDeleter interface {
	Delete(ctx context.Context, sourcePath domain.SourcePath, target domain.EntryPath) error
}
