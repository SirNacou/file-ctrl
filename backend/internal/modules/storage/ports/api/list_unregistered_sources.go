package api

import (
	"context"
	"fmt"
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
		Method:        http.MethodGet,
		Path:          "/list-unregistered-sources",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *ListUnregisteredSourcesRequest) (*ListUnregisteredSourcesResponse, error) {
		res, err := handler.Handle(ctx, &app.ListUnregisteredSourcesRequest{})
		if err != nil {
			return nil, err
		}
		fmt.Printf("res: %q\n", res.Dirs)

		return &ListUnregisteredSourcesResponse{
			ListUnregisteredSourcesResponseBody{
				Dirs: res.Dirs,
			},
		}, nil
	})
}
