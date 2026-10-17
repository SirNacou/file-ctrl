package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type ListDirectoryRequest struct {
	Path string `query:"path" default:"."`

	StorageSubpath string `path:"storage_subpath"`
}

type ListDirectoryResponse struct {
	Body app.ListDirectoryResponse
}

func ListDirectory(api huma.API, handler *app.ListDirectoryHandler) {
	huma.Register(api, huma.Operation{
		OperationID:   "list-directory",
		Method:        http.MethodGet,
		Path:          "/{storage_subpath}/files",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *ListDirectoryRequest) (*ListDirectoryResponse, error) {
		res, err := handler.Handle(ctx, &app.ListDirectoryRequest{
			StorageSubpath: i.StorageSubpath,
			Path:           i.Path,
		})
		if err != nil {
			return nil, err
		}

		return &ListDirectoryResponse{
			Body: *res,
		}, nil
	})
}
