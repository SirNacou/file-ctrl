package app

import (
	"context"
)

type ValidateUploadRequest struct {
	StoragePath string
	ParentPath  string
	Filename    string
}
type (
	ValidateUploadResponse struct{}
	ValidateUploadHandler  struct{}
)

func NewValidateUploadHandler() *ValidateUploadHandler {
	return &ValidateUploadHandler{}
}

func (h *ValidateUploadHandler) Handle(
	ctx context.Context,
	req *ValidateUploadRequest,
) (*ValidateUploadResponse, error) {
	return &ValidateUploadResponse{}, nil
}
