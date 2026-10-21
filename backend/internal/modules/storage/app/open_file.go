package app

import (
	"context"
	"io"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type (
	OpenFileRequest struct {
		StoragePath string
		FilePath    string
	}
	OpenFileResponse struct {
		File     io.ReadSeekCloser
		FileMeta domain.FileMeta
	}
	OpenFileHandler struct {
		storageSourcesFinder StorageSourcesFinder
		fileOpener           FileOpener
	}
)

func NewOpenFileHandler(
	storageSourcesFinder StorageSourcesFinder,
	fileOpener FileOpener,
) *OpenFileHandler {
	return &OpenFileHandler{
		storageSourcesFinder: storageSourcesFinder,
		fileOpener:           fileOpener,
	}
}

func (h *OpenFileHandler) Handle(
	ctx context.Context,
	req *OpenFileRequest,
) (*OpenFileResponse, error) {
	source, err := h.storageSourcesFinder.Find(ctx, req.StoragePath)
	if err != nil {
		return nil, domain.ErrStorageSourceNotFound
	}

	targetPath, err := domain.NewEntryPath(req.FilePath)
	if err != nil {
		return nil, err
	}

	file, info, err := h.fileOpener.OpenFile(ctx, source.Path, targetPath)
	if err != nil {
		return nil, err
	}

	return &OpenFileResponse{
		File:     file,
		FileMeta: info,
	}, nil
}
