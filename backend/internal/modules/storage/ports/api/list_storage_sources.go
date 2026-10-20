package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type (
	ListStorageSourcesRequest  struct{}
	ListStorageSourcesResponse struct {
		Body ListStorageSourcesResponseBody
	}
	ListStorageSourcesResponseBody struct {
		Sources []app.Source `json:"sources"`
	}
)

func ListStorageSources(api huma.API, handler *app.ListStorageSourcesHandler) {
	huma.Register(api, huma.Operation{
		OperationID: "list-storage-sources",
		Method:        http.MethodGet,
		Path:          "",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *ListStorageSourcesRequest) (*ListStorageSourcesResponse, error) {
		res, err := handler.Handle(ctx, &app.ListStorageSourcesRequest{})
		if err != nil {
			return nil, err
		}

		return &ListStorageSourcesResponse{
			Body: ListStorageSourcesResponseBody{
				Sources: res.Sources,
			},
		}, nil
	})
}
