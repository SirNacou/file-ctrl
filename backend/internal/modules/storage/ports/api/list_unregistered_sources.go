package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type ListUnregisteredSourcesRequest struct{}
type ListUnregisteredSourcesResponse struct {
	Dirs []string `json:"dirs"`
}

func ListUnregisteredSources(api huma.API, handler *app.ListUnregisteredSourcesHandler) {
	huma.Register(api, huma.Operation{
		Method:        "GET",
		Path:          "/list-unregistered-sources",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *ListUnregisteredSourcesRequest) (*ListUnregisteredSourcesResponse, error) {
		res, err := handler.Handle(ctx, &app.ListUnregisteredSourcesRequest{})
		if err != nil {
			return nil, err
		}
		return &ListUnregisteredSourcesResponse{
			Dirs: res.Dirs,
		}, nil
	})
}
