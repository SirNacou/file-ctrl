package app

import (
	"context"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type ValidateUploadRequest struct {
	StoragePath  string
	StorageSubpath   string
	DirectoryPath string
	Filename     string
}
type (
	ValidateUploadResponse struct{}
	ValidateUploadHandler  struct {
		finder StorageSourcesFinder
	}
)

func NewValidateUploadHandler(finder StorageSourcesFinder) *ValidateUploadHandler {
	return &ValidateUploadHandler{
		finder: finder,
	}
}

func (h *ValidateUploadHandler) Handle(
	ctx context.Context,
	req *ValidateUploadRequest,
) (*ValidateUploadResponse, error) {
	if _, err := h.finder.Find(ctx, req.StoragePath); err != nil {
		return nil, domain.ErrStorageSourceNotFound
	}

	sourceSubpath, err := domain.NewEntryPath(req.StorageSubpath)
	if err != nil {
		return nil, err
	}

	uploadDirectory, err := sourceSubpath.Join(req.DirectoryPath)
	if err != nil {
		return nil, err
	}

	if _, err := uploadDirectory.Join(req.Filename); err != nil {
		return nil, err
	}

	return &ValidateUploadResponse{}, nil
}
