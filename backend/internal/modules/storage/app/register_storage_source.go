package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	subpathVO, err := domain.NewSourceSubpath(req.SubPath, false)
	if err != nil {
		return nil, err
	}

	diskPath := filepath.Join(h.storageDir, filepath.FromSlash(subpathVO.String()))

	if err = h.assertDirectoryExists(diskPath); err != nil {
		return nil, err
	}

	source, err := domain.NewStorageSource(req.Name, *subpathVO)
	if err != nil {
		return nil, err
	}

	err = h.inserter.Insert(ctx, source)
	if err != nil {
		return nil, err
	}

	return &RegisterStorageSourceResponse{
		ID: source.ID,
	}, nil
}

func (h *RegisterStorageSourceHandler) assertDirectoryExists(targetDiskPath string) error {
	info, err := os.Stat(targetDiskPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("directory does not exist on disk")
		}

		return fmt.Errorf("filesystem error: %w", err)
	}

	if !info.IsDir() {
		return errors.New("target path exists but is not a directory")
	}

	// Symlink boundary verification
	realTarget, err := filepath.EvalSymlinks(targetDiskPath)
	if err != nil {
		return fmt.Errorf("cannot resolve symlink: %w", err)
	}

	realRoot, err := filepath.EvalSymlinks(h.storageDir)
	if err != nil {
		return fmt.Errorf("cannot resolve storage root symlink: %w", err)
	}

	rel, err := filepath.Rel(realRoot, realTarget)
	if err != nil || strings.HasPrefix(rel, "..") {
		return errors.New("symlink escapes the storage root")
	}

	return nil
}
