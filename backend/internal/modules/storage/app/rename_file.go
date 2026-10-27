package app

import (
	"context"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type RenameFileRequest struct {
	StoragePath string
	ParentPath  string
	OldName     string
	NewName     string
}
type (
	RenameFileResponse struct{}
	RenameFileHandler  struct {
		finder        StorageSourcesFinder
		fsFileRenamer FileRenamer
	}
)

func NewRenameFileHandler(
	finder StorageSourcesFinder,
	fsFileRenamer FileRenamer,
) *RenameFileHandler {
	return &RenameFileHandler{
		finder:        finder,
		fsFileRenamer: fsFileRenamer,
	}
}

func (h *RenameFileHandler) Handle(
	ctx context.Context,
	req *RenameFileRequest,
) (*RenameFileResponse, error) {
	source, err := h.finder.Find(ctx, req.StoragePath)
	if err != nil {
		return nil, domain.ErrStorageSourceNotFound
	}

	parent, err := domain.NewEntryPath(req.ParentPath)
	if err != nil {
		return nil, err
	}
	oldPath, err := parent.Join(req.OldName)
	if err != nil {
		return nil, err
	}
	newPath, err := parent.Join(req.NewName)
	if err != nil {
		return nil, err
	}
	if oldPath.FS() == newPath.FS() {
		return nil, domain.ErrAlreadyExists
	}

	if err = h.fsFileRenamer.Rename(ctx, source.Path, oldPath, newPath); err != nil {
		return nil, err
	}

	return &RenameFileResponse{}, nil
}
