package app

import (
	"context"
	"errors"
	"io/fs"

	"github.com/sirnacou/file-ctrl/backend/internal/common/utils"
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
		finder StorageSourcesFinder
		opener SourceOpener
	}
)

func NewRenameFileHandler(finder StorageSourcesFinder, opener SourceOpener) *RenameFileHandler {
	return &RenameFileHandler{
		finder: finder,
		opener: opener,
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
	oldPath, err := parent.ChildPath(req.OldName)
	if err != nil {
		return nil, err
	}
	newPath, err := parent.ChildPath(req.NewName)
	if err != nil {
		return nil, err
	}
	if oldPath.FS() == newPath.FS() {
		return nil, domain.ErrAlreadyExists
	}

	root, err := h.opener.OpenSource(source.Path.String())
	if err != nil {
		return nil, domain.ErrStorageSourceNotFound
	}
	defer utils.CloseLog(ctx, root, "source root")

	// Source must exist.
	if _, err := root.Lstat(oldPath.FS()); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, domain.ErrNotFound
		}
		return nil, domain.ErrInvalidPath
	}

	// Rename silently overwrites files on Unix, so refuse if the target exists.
	if _, err := root.Lstat(newPath.FS()); err == nil {
		return nil, domain.ErrAlreadyExists
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, domain.ErrInvalidPath
	}

	if err := root.Rename(oldPath.FS(), newPath.FS()); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return nil, domain.ErrAlreadyExists
		}
		return nil, domain.ErrInvalidPath
	}

	return &RenameFileResponse{}, nil
}
