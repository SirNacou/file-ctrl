package app

import (
	"context"

	"github.com/sirnacou/file-ctrl/backend/internal/common/localfs"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type RegisterStorageSourceHandler struct {
	inserter     StorageSourcesInserter
	sourceOpener SourceOpener
}

type RegisterStorageSourceRequest struct {
	Name string
	Path string
}

type RegisterStorageSourceResponse struct {
	Path domain.StoragePath
}

func NewRegisterStorageSourceHandler(
	inserter StorageSourcesInserter,
	sourceOpener SourceOpener,
) *RegisterStorageSourceHandler {
	return &RegisterStorageSourceHandler{
		inserter:     inserter,
		sourceOpener: sourceOpener,
	}
}

func (h *RegisterStorageSourceHandler) Handle(
	ctx context.Context,
	req *RegisterStorageSourceRequest,
) (*RegisterStorageSourceResponse, error) {
	path, err := domain.NewSourcePath(req.Path, false)
	if err != nil {
		return nil, err
	}

	root, err := h.sourceOpener.OpenSource(path.String())
	if err != nil {
		return nil, domain.ErrInvalidSourceDirectory
	}
	localfs.CloseLog(ctx, root, "source root")

	source, err := domain.NewStorageSource(path, req.Name)
	if err != nil {
		return nil, err
	}
	if err = h.inserter.Insert(ctx, source); err != nil {
		return nil, err
	}

	return &RegisterStorageSourceResponse{Path: source.Path}, nil
}
