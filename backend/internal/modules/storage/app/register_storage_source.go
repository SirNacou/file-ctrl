package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"uuid"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type StorageSourcesInserter interface {
	Insert(ctx context.Context, model *domain.StorageSource) error
}

type RegisterStorageSourceHandler struct {
	inserter   StorageSourcesInserter
	storageDir string
}

type RegisterStorageSourceRequest struct {
	Name    string
	SubPath string
}

type RegisterStorageSourceResponse struct {
	ID uuid.UUID
}

func NewRegisterStorageSourceHandler(
	inserter StorageSourcesInserter,
	storageDir string,
) *RegisterStorageSourceHandler {
	return &RegisterStorageSourceHandler{
		inserter:   inserter,
		storageDir: storageDir,
	}
}

func (h *RegisterStorageSourceHandler) Handle(
	ctx context.Context,
	req *RegisterStorageSourceRequest,
) (*RegisterStorageSourceResponse, error) {
	subpath, err := domain.NewSourceSubpath(req.SubPath, false)
	if err != nil {
		return nil, err
	}

	if err = h.assertDirectoryExists(ctx, subpath.String()); err != nil {
		return nil, err
	}

	source, err := domain.NewStorageSource(req.Name, subpath)
	if err != nil {
		return nil, err
	}

	if err = h.inserter.Insert(ctx, source); err != nil {
		return nil, err
	}

	return &RegisterStorageSourceResponse{ID: source.ID}, nil
}

// assertDirectoryExists checks that subPath (relative, already validated)
// is an existing directory inside storageDir. os.Root rejects "..", absolute
// paths and symlinks that point outside the root.
func (h *RegisterStorageSourceHandler) assertDirectoryExists(
	ctx context.Context,
	subPath string,
) error {
	root, err := os.OpenRoot(h.storageDir)
	if err != nil {
		return fmt.Errorf("cannot open storage root: %w", err)
	}
	defer func() {
		err = root.Close()
		if err != nil {
			slog.ErrorContext(ctx, "failed to close root dir", slog.Any("error", err))
		}
	}()

	info, err := root.Stat(filepath.FromSlash(subPath))
	if err != nil || !info.IsDir() {
		return domain.ErrInvalidSourceDirectory
	}

	return nil
}
