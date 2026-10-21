package api

import (
	"context"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type DownloadFileRequest struct {
	StoragePath string `path:"storage_path" doc:"Storage source identifier"`

	FilePath string `query:"file_path" required:"true" doc:"Relative file path to download"`
}

func DownloadFile(api huma.API, handler *app.DownloadFileHandler) {
	huma.Register(api, huma.Operation{
		OperationID:   "storage-file-download",
		Method:        http.MethodGet,
		Path:          "/{storage_path}/files/download",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *DownloadFileRequest) (*huma.StreamResponse, error) {
		res, err := handler.Handle(ctx, &app.DownloadFileRequest{
			StoragePath: i.StoragePath,
			FilePath:    i.FilePath,
		})
		if err != nil {
			return nil, err
		}

		return &huma.StreamResponse{
			Body: func(ctx huma.Context) {
				defer res.File.Close()

				e := humaecho.Unwrap(ctx)
				h := e.Response().Header()

				h.Set(echo.HeaderXContentTypeOptions, "nosniff")
				h.Set(
					echo.HeaderContentDisposition,
					mime.FormatMediaType("attachment", map[string]string{
						"filename": res.FileEntry.Name,
					}),
				)

				contentType := mime.TypeByExtension(filepath.Ext(res.FileEntry.Name))
				if contentType == "" {
					contentType = echo.MIMEOctetStream
				}
				h.Set(echo.HeaderContentType, contentType)

				http.ServeContent(
					e.Response(),
					e.Request(),
					res.FileEntry.Name,
					res.FileEntry.ModTime,
					res.File,
				)
			},
		}, nil
	})
}
