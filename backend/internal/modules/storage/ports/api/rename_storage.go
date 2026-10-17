package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type (
	RenameStorageRequest struct {
		Subpath string `path:"subpath"`
		Body    RenameStorageRequestBody
	}
	RenameStorageResponse struct{}

	RenameStorageRequestBody struct {
		NewName string `json:"new_name"`
	}
)

func RenameStorage(api huma.API, handler *app.RenameStorageHandler) {
	huma.Register(api, huma.Operation{
		OperationID:   "storage-rename",
		Method:        http.MethodPost,
		Path:          "/{subpath}/rename",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *RenameStorageRequest) (*RenameStorageResponse, error) {
		_, err := handler.Handle(ctx, &app.RenameStorageRequest{
			Subpath: i.Subpath,
			NewName: i.Body.NewName,
		})

		return nil, err
	})
}
