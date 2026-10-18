package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type RegisterStorageSourceRequest struct {
	Body struct {
		Name string `json:"name"`
		Path string `json:"path"`
	}
}

type RegisterStorageSourceResponse struct {
	Body RegisterStorageSourceResponseBody
}

type RegisterStorageSourceResponseBody struct {
	Path string `json:"path"`
}

func RegisterStorageSource(api huma.API, handler *app.RegisterStorageSourceHandler) {
	huma.Register(api, huma.Operation{
		OperationID:   "register-storage",
		Method:        http.MethodPost,
		Path:          "",
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, i *RegisterStorageSourceRequest) (*RegisterStorageSourceResponse, error) {
		res, err := handler.Handle(ctx, &app.RegisterStorageSourceRequest{
			Name: i.Body.Name,
			Path: i.Body.Path,
		})
		if err != nil {
			return nil, err
		}

		return &RegisterStorageSourceResponse{
			Body: RegisterStorageSourceResponseBody{
				Path: res.Path.String(),
			},
		}, nil
	})
}
