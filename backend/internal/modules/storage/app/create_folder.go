package app

import (
	"context"
	"errors"
	"io/fs"

	"github.com/sirnacou/file-ctrl/backend/internal/common/localfs"
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
		finder StorageSourcesFinder
		opener SourceOpener
	}
)

func NewCreateFolderHandler(finder StorageSourcesFinder, opener SourceOpener) *CreateFolderHandler {
	return &CreateFolderHandler{
		finder: finder,
		opener: opener,
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

	root, err := h.opener.OpenSource(source.Path.String())
	if err != nil {
		return nil, domain.ErrStorageSourceNotFound
	}
	defer localfs.CloseLog(ctx, root, "source root")

	if err := root.Mkdir(target.FS(), 0o755); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return nil, domain.ErrAlreadyExists
		}
		return nil, domain.ErrInvalidPath
	}

	return &CreateFolderResponse{}, nil
}
