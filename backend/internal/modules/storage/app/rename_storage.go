package app

import (
	"context"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type (
	RenameStorageRequest struct {
		Path    string
		NewName string
	}
	RenameStorageResponse struct{}
	RenameStorageHandler  struct {
		finder  StorageSourcesFinder
		updater StorageSourcesUpdater
	}
)

func NewRenameStorageHandler(
	finder StorageSourcesFinder,
	updater StorageSourcesUpdater,
) *RenameStorageHandler {
	return &RenameStorageHandler{
		finder:  finder,
		updater: updater,
	}
}

func (h *RenameStorageHandler) Handle(
	ctx context.Context,
	req *RenameStorageRequest,
) (*RenameStorageResponse, error) {
	storage, err := h.finder.Find(ctx, req.Path)
	if err != nil {
		return nil, domain.ErrStorageSourceNotFound
	}

	err = storage.Rename(req.NewName)
	if err != nil {
		return nil, err
	}

	err = h.updater.Update(ctx, storage)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
