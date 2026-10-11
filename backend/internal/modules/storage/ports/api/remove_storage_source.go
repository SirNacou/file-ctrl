package api

import (
	"context"
	"net/http"
	"uuid"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type RemoveStorageSourceRequest struct {
	ID uuid.UUID `json:"id" format:"uuid"`
}
type RemoveStorageSourceResponse struct{}

func RemoveStorageSource(api huma.API, handler *app.RemoveStorageSourceHandler) {
	huma.Register(api, huma.Operation{
		Method:        http.MethodDelete,
		Path:          "/remove-storage-source",
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, i *RemoveStorageSourceRequest) (*RemoveStorageSourceResponse, error) {
		_, err := handler.Handle(ctx, &app.RemoveStorageSourceRequest{
			ID: i.ID,
		})
		if err != nil {
			return nil, err
		}

		return &RemoveStorageSourceResponse{}, nil
	})
}
