package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type CreateFolderRequest struct {
	StoragePath string `path:"storage_path"`
	Body        struct {
		ParentPath string `json:"parent_path"`
		Name       string `json:"name"`
	}
}
type CreateFolderResponse struct{}

func CreateFolder(api huma.API, handler *app.CreateFolderHandler) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-storage-folder",
		Method:        http.MethodPost,
		Path:          "/folder",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *CreateFolderRequest) (*CreateFolderResponse, error) {
		_, err := handler.Handle(ctx, &app.CreateFolderRequest{
			Name:        i.Body.Name,
			StoragePath: i.StoragePath,
			ParentPath:  i.Body.ParentPath,
		})

		return nil, err
	})
}
