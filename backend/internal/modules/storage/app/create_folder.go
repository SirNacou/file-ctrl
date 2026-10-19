package app

import (
	"context"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type CreateFolderRequest struct {
	Name        string
	StoragePath string
	ParentPath  string
}
type (
	CreateFolderResponse struct{}
	CreateFolderHandler  struct {
		finder       StorageSourcesFinder
		fsDirCreator FSDirCreator
	}
)

func NewCreateFolderHandler(
	finder StorageSourcesFinder,
	fsDirCreator FSDirCreator,
) *CreateFolderHandler {
	return &CreateFolderHandler{
		finder:       finder,
		fsDirCreator: fsDirCreator,
	}
}

func (h *CreateFolderHandler) Handle(
	ctx context.Context,
	req *CreateFolderRequest,
) (*CreateFolderResponse, error) {
	source, err := h.finder.Find(ctx, req.StoragePath)
	if err != nil {
		return nil, domain.ErrStorageSourceNotFound
	}

	parent, err := domain.NewEntryPath(req.ParentPath)
	if err != nil {
		return nil, err
	}
	target, err := parent.ChildPath(req.Name)
	if err != nil {
		return nil, err
	}

	err = h.fsDirCreator.Mkdir(ctx, source.Path, target)
	if err != nil {
		return nil, err
	}

	return &CreateFolderResponse{}, nil
}
