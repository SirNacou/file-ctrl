package app

import (
	"context"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type RegisterStorageSourceHandler struct {
	inserter           StorageSourcesInserter
	fsDirExistsChecker DirExistsChecker
}

type RegisterStorageSourceRequest struct {
	Name string
	Path string
}

type RegisterStorageSourceResponse struct {
	Path domain.SourcePath
}

func NewRegisterStorageSourceHandler(
	inserter StorageSourcesInserter,
	fsDirExistsChecker DirExistsChecker,
) *RegisterStorageSourceHandler {
	return &RegisterStorageSourceHandler{
		inserter:           inserter,
		fsDirExistsChecker: fsDirExistsChecker,
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

	exists, err := h.fsDirExistsChecker.DirExists(ctx, path.String())
	if err != nil {
		return nil, domain.ErrInvalidSourceDirectory
	}
	if !exists {
		return nil, domain.ErrSourceDirectoryNotFound
	}

	source, err := domain.NewStorageSource(path, req.Name)
	if err != nil {
		return nil, err
	}
	if err = h.inserter.Insert(ctx, source); err != nil {
		return nil, err
	}

	return &RegisterStorageSourceResponse{Path: source.Path}, nil
}
