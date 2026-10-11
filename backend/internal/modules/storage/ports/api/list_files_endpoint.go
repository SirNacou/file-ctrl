package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type Request struct {
	Path string `query:"path"`
}

type Response struct {
	Body app.ListFilesResponse
}

func ListFile(api huma.API, handler *app.ListFilesHandler) {
	huma.Register(api, huma.Operation{
		Method:        "GET",
		Path:          "/list-files",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *Request) (*Response, error) {

		res, err := handler.Handle(ctx, &app.ListFilesRequest{Path: i.Path})
		if err != nil {
			return nil, err
		}

		return &Response{
			Body: *res,
		}, nil
	})
}
