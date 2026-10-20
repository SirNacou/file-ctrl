package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type (
	ListUnregisteredSourcesRequest  struct{}
	ListUnregisteredSourcesResponse struct {
		Body ListUnregisteredSourcesResponseBody
	}
	ListUnregisteredSourcesResponseBody struct {
		Dirs []string `json:"dirs"`
	}
)

func ListUnregisteredSources(api huma.API, handler *app.ListUnregisteredSourcesHandler) {
	huma.Register(api, huma.Operation{
		OperationID:   "list-unregistered-sources",
		Method:        http.MethodGet,
		Path:          "/unregistered-sources",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *ListUnregisteredSourcesRequest) (*ListUnregisteredSourcesResponse, error) {
		res, err := handler.Handle(ctx, &app.ListUnregisteredSourcesRequest{})
		if err != nil {
			return nil, err
		}

		return &ListUnregisteredSourcesResponse{
			ListUnregisteredSourcesResponseBody{
				Dirs: res.Dirs,
			},
		}, nil
	})
}
