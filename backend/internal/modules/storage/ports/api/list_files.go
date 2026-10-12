package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type ListFileRequest struct {
	Path string `query:"path"`
}

type ListFileResponse struct {
	Body app.ListFilesResponse
}

func ListFile(api huma.API, handler *app.ListFilesHandler) {
	huma.Register(api, huma.Operation{
		Method:        http.MethodGet,
		Path:          "/list-files",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *ListFileRequest) (*ListFileResponse, error) {
		res, err := handler.Handle(ctx, &app.ListFilesRequest{Path: i.Path})
		if err != nil {
			return nil, err
		}

		return &ListFileResponse{
			Body: *res,
		}, nil
	})
}
