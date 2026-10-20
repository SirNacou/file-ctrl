package app

import (
	"context"
	"io"
	"os"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type (
	DownloadFileRequest struct {
		StoragePath string
		FilePath    string
	}
	DownloadFileResponse struct {
		File     io.ReadCloser
		FileInfo os.FileInfo
		Filename string
	}
	DownloadFileHandler struct {
		storageSourcesFinder StorageSourcesFinder
		fileOpener           FileOpener
	}
)

func NewDownloadFileHandler(
	storageSourcesFinder StorageSourcesFinder,
	fileOpener FileOpener,
) *DownloadFileHandler {
	return &DownloadFileHandler{
		storageSourcesFinder: storageSourcesFinder,
		fileOpener:           fileOpener,
	}
}

func (h *DownloadFileHandler) Handle(
	ctx context.Context,
	req *DownloadFileRequest,
) (*DownloadFileResponse, error) {
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

	return &DownloadFileResponse{
		File:     file,
		FileInfo: info,
		Filename: file.Name(),
	}, nil
}
