package api

import (
	"context"
	"net/http"
	"uuid"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type RegisterStorageSourceRequest struct {
	Name    string `json:"name"`
	SubPath string `json:"sub_path"`
}
type RegisterStorageSourceResponse struct {
	ID uuid.UUID `json:"id"`
}

func RegisterStorageSource(api huma.API, handler *app.RegisterStorageSourceHandler) {
	huma.Register(api, huma.Operation{
		Method:        http.MethodPost,
		Path:          "",
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, i *RegisterStorageSourceRequest) (*RegisterStorageSourceResponse, error) {
		res, err := handler.Handle(ctx, &app.RegisterStorageSourceRequest{
			Name:    i.Name,
			SubPath: i.SubPath,
		})
		if err != nil {
			return nil, err
		}

		return &RegisterStorageSourceResponse{
			ID: res.ID,
		}, nil
	})
}
