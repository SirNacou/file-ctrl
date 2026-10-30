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
		finderupdater StorageSourcesFinderUpdater
	}
)

func NewRenameStorageHandler(
	finderupdater StorageSourcesFinderUpdater,
) *RenameStorageHandler {
	return &RenameStorageHandler{
		finderupdater: finderupdater,
	}
}

func (h *RenameStorageHandler) Handle(
	ctx context.Context,
	req *RenameStorageRequest,
) (*RenameStorageResponse, error) {
	storage, err := h.finderupdater.Find(ctx, req.Path)
	if err != nil {
		return nil, domain.ErrStorageSourceNotFound
	}

	err = storage.Rename(req.NewName)
	if err != nil {
		return nil, err
	}

	err = h.finderupdater.Update(ctx, storage)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
