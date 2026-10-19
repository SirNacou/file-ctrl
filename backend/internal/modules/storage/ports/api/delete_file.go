package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type DeleteFileRequest struct {
	StoragePath string `path:"storage_path"`
	Body        struct {
		FilePath string `json:"file_path"`
	}
}
type DeleteFileResponse struct{}

func DeleteFile(api huma.API, handler *app.DeleteFileHandler) {
	huma.Register(api, huma.Operation{
		OperationID:   "storage-delete-file",
		Method:        http.MethodDelete,
		DefaultStatus: http.StatusAccepted,
		Path:          "/{storage_path}/files",
	}, func(ctx context.Context, i *DeleteFileRequest) (*DeleteFileResponse, error) {
		_, err := handler.Handle(ctx, &app.DeleteFileRequest{
			StoragePath: i.StoragePath,
			FilePath:    i.Body.FilePath,
		})

		return &DeleteFileResponse{}, err
	})
}
