package app

import (
	"context"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
	"github.com/stephenafamo/bob"
)

type (
	RenameStorageRequest struct {
		Subpath string
		NewName string
	}
	RenameStorageResponse struct{}
	RenameStorageHandler  struct {
		db                bob.Executor
		storageSourceRepo StorageSourceRepo
		storageDir        string
	}
	StorageSourceRepo interface {
		Find(
			ctx context.Context,
			subpath string,
		) (*domain.StorageSource, error)
		Update(ctx context.Context, model *domain.StorageSource) error
	}
)

func NewRenameStorageHandler(
	repo StorageSourceRepo,
	db bob.Executor,
	storageDir string,
) *RenameStorageHandler {
	return &RenameStorageHandler{
		storageSourceRepo: repo,
		db:                db,
		storageDir:        storageDir,
	}
}

func (h *RenameStorageHandler) Handle(
	ctx context.Context,
	req *RenameStorageRequest,
) (*RenameStorageResponse, error) {
	storage, err := h.storageSourceRepo.Find(ctx, req.Subpath)
	if err != nil {
		return nil, domain.ErrStorageSourceNotFound
	}

	err = storage.Rename(req.NewName)
	if err != nil {
		return nil, err
	}

	err = h.storageSourceRepo.Update(ctx, storage)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
