package app

import (
	"context"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type DeleteFilesRequest struct {
	StoragePath string
	FilePaths   []string
}

type (
	DeleteFilesResponse struct{}

	DeleteFilesHandler struct {
		finder        StorageSourcesFinder
		fsFileDeleter FileDeleter
	}
)

func NewDeleteFilesHandler(
	finder StorageSourcesFinder,
	fsFileDeleter FileDeleter,
) *DeleteFilesHandler {
	return &DeleteFilesHandler{
		finder:        finder,
		fsFileDeleter: fsFileDeleter,
	}
}

func (h *DeleteFilesHandler) Handle(
	ctx context.Context,
	req *DeleteFilesRequest,
) (*DeleteFilesResponse, error) {
	source, err := h.finder.Find(ctx, req.StoragePath)
	if err != nil {
		return nil, domain.ErrStorageSourceNotFound
	}

	for _, path := range req.FilePaths {
		entryPath, err := domain.NewEntryPath(path)
		if err != nil {
			return nil, err
		}

		if err := h.fsFileDeleter.Delete(ctx, source.Path, entryPath); err != nil {
			return nil, err
		}
	}

	return &DeleteFilesResponse{}, nil
}
