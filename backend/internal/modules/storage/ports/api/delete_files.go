package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type DeleteFilesRequest struct {
	StoragePath string `path:"storage_path"`
	Body        struct {
		FilePaths []string `json:"file_paths" minLength:"1"`
	}
}
type DeleteFilesResponse struct{}

func DeleteFiles(api huma.API, handler *app.DeleteFilesHandler) {
	huma.Register(api, huma.Operation{
		OperationID:   "storage-delete-files",
		Method:        http.MethodDelete,
		DefaultStatus: http.StatusAccepted,
		Path:          "/{storage_path}/files",
	}, func(ctx context.Context, i *DeleteFilesRequest) (*DeleteFilesResponse, error) {
		_, err := handler.Handle(ctx, &app.DeleteFilesRequest{
			StoragePath: i.StoragePath,
			FilePaths:   i.Body.FilePaths,
		})

		return &DeleteFilesResponse{}, err
	})
}
