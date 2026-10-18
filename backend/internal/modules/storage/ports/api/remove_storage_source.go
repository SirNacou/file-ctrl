package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type RemoveStorageSourceRequest struct {
	Path string `path:"path"`
}
type RemoveStorageSourceResponse struct{}

func RemoveStorageSource(api huma.API, handler *app.RemoveStorageSourceHandler) {
	huma.Register(api, huma.Operation{
		OperationID:   "remove-storage",
		Method:        http.MethodDelete,
		Path:          "/{path}",
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, i *RemoveStorageSourceRequest) (*RemoveStorageSourceResponse, error) {
		_, err := handler.Handle(ctx, &app.RemoveStorageSourceRequest{
			Path: i.Path,
		})
		if err != nil {
			return nil, err
		}

		return &RemoveStorageSourceResponse{}, nil
	})
}
