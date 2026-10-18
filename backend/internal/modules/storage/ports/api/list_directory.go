package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type ListDirectoryRequest struct {
	Subpath string `query:"subpath" default:"."`

	StoragePath string `path:"storage_path"`
}

type ListDirectoryResponse struct {
	Body app.ListDirectoryResponse
}

func ListDirectory(api huma.API, handler *app.ListDirectoryHandler) {
	huma.Register(api, huma.Operation{
		OperationID:   "list-directory",
		Method:        http.MethodGet,
		Path:          "/{storage_path}/files",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *ListDirectoryRequest) (*ListDirectoryResponse, error) {
		res, err := handler.Handle(ctx, &app.ListDirectoryRequest{
			StoragePath: i.StoragePath,
			Subpath:     i.Subpath,
		})
		if err != nil {
			return nil, err
		}

		return &ListDirectoryResponse{
			Body: *res,
		}, nil
	})
}
