package app

import (
	"context"
)

type (
	UploadFileRequest struct {
		StoragePath    string
		ParentPath     string
		Filename       string
		StagedFilePath string
	}
	UploadFileResponse struct{}
	UploadFileHandler  struct{}
)

func NewUploadFileHandler() *UploadFileHandler {
	return &UploadFileHandler{}
}

func (h *UploadFileHandler) Handle(
	ctx context.Context,
	req *UploadFileRequest,
) (*UploadFileResponse, error) {
	return &UploadFileResponse{}, nil
}
