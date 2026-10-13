package api

import (
	"context"
	"net/http"
	"uuid"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type ListDirectoryRequest struct {
	Path string    `query:"path" default:"."`
	ID   uuid.UUID `                         path:"id" format:"uuid"`
}

type ListDirectoryResponse struct {
	Body app.ListDirectoryResponse
}

func ListDirectory(api huma.API, handler *app.ListDirectoryHandler) {
	huma.Register(api, huma.Operation{
		Method:        http.MethodGet,
		Path:          "/{id}/files",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *ListDirectoryRequest) (*ListDirectoryResponse, error) {
		res, err := handler.Handle(ctx, &app.ListDirectoryRequest{
			StorageID: i.ID,
			Path:      i.Path,
		})
		if err != nil {
			return nil, err
		}

		return &ListDirectoryResponse{
			Body: *res,
		}, nil
	})
}
