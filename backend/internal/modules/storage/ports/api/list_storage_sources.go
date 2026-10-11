package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type ListStorageSourcesRequest struct{}
type ListStorageSourcesResponse struct {
	Sources []app.Source `json:"sources"`
}

func ListStorageSources(api huma.API, handler *app.ListStorageSourcesHandler) {
	huma.Register(api, huma.Operation{
		Method:        http.MethodGet,
		Path:          "",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *ListStorageSourcesRequest) (*ListStorageSourcesResponse, error) {
		res, err := handler.Handle(ctx, &app.ListStorageSourcesRequest{})
		if err != nil {
			return nil, err
		}

		return &ListStorageSourcesResponse{
			Sources: res.Sources,
		}, nil
	})
}
