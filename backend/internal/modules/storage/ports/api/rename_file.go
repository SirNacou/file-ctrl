package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type RenameFileRequest struct {
	StoragePath string `path:"storage_path"`
	Body        struct {
		ParentPath string `json:"parent_path"`
		OldName    string `json:"old_name"`
		NewName    string `json:"new_name"`
	}
}
type RenameFileResponse struct{}

func RenameFile(api huma.API, handler *app.RenameFileHandler) {
	huma.Register(api, huma.Operation{
		OperationID:   "rename-file",
		Method:        http.MethodPost,
		Path:          "/{storage_path}/files/rename",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *RenameFileRequest) (*RenameFileResponse, error) {
		_, err := handler.Handle(ctx, &app.RenameFileRequest{
			StoragePath: i.StoragePath,
			ParentPath:  i.Body.ParentPath,
			OldName:     i.Body.OldName,
			NewName:     i.Body.NewName,
		})

		return nil, err
	})
}
