package api

import (
	"context"
	"log/slog"
	"mime"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

type DownloadArchiveRequest struct {
	StoragePath string `path:"storage_path"`

	FilePaths []string `query:"file_paths" minLength:"1"`
}
type DownloadArchiveResponse struct{}

func DownloadArchive(api huma.API, handler *app.PrepareArchiveHandler) {
	huma.Register(api, huma.Operation{
		OperationID: "storage-archive",
		Method:      http.MethodGet,
		Path:        "/{storage_path}/files/archive",
		Summary:     "Download files and folders as a zip",
	}, func(ctx context.Context, i *DownloadArchiveRequest) (*huma.StreamResponse, error) {
		res, err := handler.Handle(ctx, &app.PrepareArchiveRequest{
			StoragePath: i.StoragePath,
			FilePaths:   i.FilePaths,
		})
		if err != nil {
			return nil, err
		}

		return &huma.StreamResponse{
			Body: func(hctx huma.Context) {
				defer res.Archive.Close()

				rc := http.NewResponseController(humaecho.Unwrap(hctx).Response())
				if err := rc.SetWriteDeadline(time.Time{}); err != nil {
					slog.WarnContext(
						hctx.Context(),
						"cannot clear write deadline",
						slog.Any("err", err),
					)
				}

				hctx.SetHeader("Content-Type", "application/zip")
				hctx.SetHeader("Content-Disposition", contentDisposition(res.Archive.Name()))
				hctx.SetHeader("X-Accel-Buffering", "no") // tells nginx not to buffer
				if err := res.Archive.WriteZip(hctx.Context(), hctx.BodyWriter()); err != nil {
					slog.ErrorContext(hctx.Context(), "zip stream failed", slog.Any("err", err))
				}
			},
		}, nil
	})
}

func contentDisposition(name string) string {
	// Handles quoting and RFC 6266/2231 encoding for non-ASCII names.
	return mime.FormatMediaType("attachment", map[string]string{"filename": name})
}
