package api

import (
	"context"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"

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
				header := e.Response().Header()

				header.Set(echo.HeaderXContentTypeOptions, "nosniff")
				header.Set(
					echo.HeaderContentDisposition,
					mime.FormatMediaType("attachment", map[string]string{
						"filename": res.Filename,
					}),
				)
				header.Set(
					echo.HeaderContentLength,
					strconv.FormatInt(res.FileInfo.Size(), 10),
				)

				contentType := mime.TypeByExtension(
					filepath.Ext(res.Filename),
				)
				if contentType == "" {
					contentType = echo.MIMEOctetStream
				}

				header.Set(echo.HeaderContentType, contentType)

				_, err = io.Copy(e.Response(), res.File)
				if err != nil {
					slog.ErrorContext(
						ctx.Context(),
						"failed to download file",
						slog.Any("error", err),
					)
				}
			},
		}, nil
	})
}
