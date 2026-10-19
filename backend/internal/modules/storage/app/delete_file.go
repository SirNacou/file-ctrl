package app

import (
	"context"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type DeleteFileRequest struct {
	StoragePath string
	FilePath    string
}

type (
	DeleteFileResponse struct{}
	DeleteFileHandler  struct {
		finder        StorageSourcesFinder
		fsFileDeleter FSFileDeleter
	}
)

func NewDeleteFileHandler(
	finder StorageSourcesFinder,
	fsFileDeleter FSFileDeleter,
) *DeleteFileHandler {
	return &DeleteFileHandler{
		finder:        finder,
		fsFileDeleter: fsFileDeleter,
	}
}

func (h *DeleteFileHandler) Handle(
	ctx context.Context,
	req *DeleteFileRequest,
) (*DeleteFileResponse, error) {
	source, err := h.finder.Find(ctx, req.StoragePath)
	if err != nil {
		return nil, domain.ErrStorageSourceNotFound
	}

	filePath, err := domain.NewEntryPath(req.FilePath)
	if err != nil {
		return nil, err
	}

	if err = h.fsFileDeleter.Delete(ctx, source.Path, filePath); err != nil {
		return nil, err
	}

	return &DeleteFileResponse{}, nil
}
