package app

import (
	"context"
	"os"

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

type StorageSourcesFinderUpdater interface {
	StorageSourcesFinder
	StorageSourcesUpdater
}

type DirReader interface {
	ReadDir(
		ctx context.Context,
		sourcePath domain.SourcePath,
		dir domain.EntryPath,
	) ([]domain.FileMeta, error)
}

type DirExistsChecker interface {
	DirExists(ctx context.Context, relPath string) (bool, error)
}

type DirCreator interface {
	Mkdir(ctx context.Context, sourcePath domain.SourcePath, target domain.EntryPath) error
}

type FileRenamer interface {
	Rename(ctx context.Context, sourcePath domain.SourcePath, from, to domain.EntryPath) error
}

type DirLister interface {
	ListDirs() ([]string, error)
}

type FileDeleter interface {
	Delete(ctx context.Context, sourcePath domain.SourcePath, target domain.EntryPath) error
}

type FileOpener interface {
	OpenFile(
		ctx context.Context,
		sourcePath domain.SourcePath,
		target domain.EntryPath,
	) (*os.File, domain.FileMeta, error)
}
