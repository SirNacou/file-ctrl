package api

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
	"github.com/tus/tusd/v2/pkg/filestore"
	"github.com/tus/tusd/v2/pkg/handler"
	"golang.org/x/exp/slog"
)

const UploadBasePath = "/api/storage/uploads/"

type UploadRoute struct {
	prefix  string
	handler http.Handler
}

func NewUploadRoute(
	storageDir string,
	validateUploadHandler *app.ValidateUploadHandler,
	uploadFileHandler *app.UploadFileHandler,
) (*UploadRoute, error) {
	stagingDir := filepath.Join(storageDir, ".uploads")
	if err := os.MkdirAll(stagingDir, 0o755); err != nil {
		return nil, err
	}

	store := filestore.New(stagingDir)
	composer := handler.NewStoreComposer()
	store.UseIn(composer)

	tusHandler, err := handler.NewHandler(handler.Config{
		StoreComposer: composer,
		BasePath:      UploadBasePath,
		Logger:        slog.Default(),
		PreUploadCreateCallback: func(hook handler.HookEvent) (handler.HTTPResponse, handler.FileInfoChanges, error) {
			meta, err := mapToMetadata(hook.Upload.MetaData)
			if err != nil {
				return handler.HTTPResponse{
					StatusCode: http.StatusBadRequest,
				}, handler.FileInfoChanges{}, err
			}

			_, err = validateUploadHandler.Handle(hook.Context, &app.ValidateUploadRequest{
				StoragePath: meta.StoragePath,
				ParentPath:  meta.ParentPath,
				Filename:    meta.Filename,
			})
			if err != nil {
				slog.ErrorContext(hook.Context, "failed to validate file", slog.Any("error", err))
				status := http.StatusBadRequest
				if errors.Is(err, domain.ErrStorageSourceNotFound) {
					status = http.StatusNotFound
				}
				return handler.HTTPResponse{
					StatusCode: status,
				}, handler.FileInfoChanges{}, err
			}

			return handler.HTTPResponse{}, handler.FileInfoChanges{}, nil
		},
		PreFinishResponseCallback: func(hook handler.HookEvent) (handler.HTTPResponse, error) {
			meta, err := mapToMetadata(hook.Upload.MetaData)
			if err != nil {
				return handler.HTTPResponse{
					StatusCode: http.StatusBadRequest,
				}, err
			}

			stagedPath := hook.Upload.Storage[filestore.StorageKeyPath]
			infoPath := hook.Upload.Storage[filestore.StorageKeyInfoPath]

			defer func() {
				_ = os.Remove(stagedPath)
				_ = os.Remove(infoPath)
			}()

			_, err = uploadFileHandler.Handle(hook.Context, &app.UploadFileRequest{
				StoragePath:    meta.StoragePath,
				ParentPath:     meta.ParentPath,
				Filename:       meta.Filename,
				StagedFilePath: stagedPath,
			})
			if err != nil {
				slog.ErrorContext(hook.Context, "failed to upload file", slog.Any("error", err))
				return handler.HTTPResponse{StatusCode: http.StatusInternalServerError}, err
			}

			return handler.HTTPResponse{}, nil
		},
	})
	if err != nil {
		return nil, err
	}

	return &UploadRoute{
		prefix:  UploadBasePath,
		handler: http.StripPrefix(UploadBasePath, tusHandler),
	}, nil
}

type metadata struct {
	StoragePath string
	ParentPath  string
	Filename    string
}

func mapToMetadata(tusMeta handler.MetaData) (metadata, error) {
	res := metadata{}

	storagePath, exists := tusMeta["storage_path"]
	storagePath = strings.TrimSpace(storagePath)
	if !exists || storagePath == "" {
		return metadata{}, handler.NewError(
			"INVALID_METADATA",
			"invalid metadata: storage_path",
			http.StatusBadRequest,
		)
	}
	res.StoragePath = storagePath

	parentPath, exists := tusMeta["parent_path"]
	parentPath = strings.TrimSpace(parentPath)
	if !exists || parentPath == "" {
		return metadata{}, handler.NewError(
			"INVALID_METADATA",
			"invalid metadata: parent_path",
			http.StatusBadRequest,
		)
	}
	res.ParentPath = parentPath

	filename, exists := tusMeta["filename"]
	filename = strings.TrimSpace(filename)
	if !exists || filename == "" {
		return metadata{}, handler.NewError(
			"INVALID_METADATA",
			"invalid metadata: filename",
			http.StatusBadRequest,
		)
	}
	res.StoragePath = filename

	return res, nil
}

func (r *UploadRoute) Prefix() string        { return r.prefix }
func (r *UploadRoute) Handler() http.Handler { return r.handler }
