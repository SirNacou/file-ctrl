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

const UploadBasePath = "/api/storage/files/uploads"

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
			slog.InfoContext(hook.Context, "Metadata", slog.Any("meta", hook.Upload.MetaData))
			meta, err := mapToMetadata(hook.Upload.MetaData)
			if err != nil {
				return handler.HTTPResponse{
					StatusCode: http.StatusBadRequest,
				}, handler.FileInfoChanges{}, err
			}

			_, err = validateUploadHandler.Handle(hook.Context, &app.ValidateUploadRequest{
				StoragePath:    meta.StoragePath,
				StorageSubpath: meta.StorageSubpath,
				DirectoryPath:  meta.DirectoryPath,
				Filename:       meta.Filename,
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
			slog.InfoContext(hook.Context, "Metadata", slog.Any("meta", hook.Upload.MetaData))
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
				StorageSubpath: meta.StorageSubpath,
				DirectoryPath:  meta.DirectoryPath,
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
	StoragePath    string
	StorageSubpath string
	DirectoryPath  string
	Filename       string
}

func mapToMetadata(tusMeta handler.MetaData) (metadata, error) {
	res := metadata{}

	storagePath, exists := tusMeta["storage_path"]
	storagePath = strings.TrimSpace(storagePath)
	if !exists || storagePath == "" {
		return metadata{}, handler.NewError(
			"INVALID_METADATA",
			"metadata storage_path must not empty",
			http.StatusBadRequest,
		)
	}
	res.StoragePath = storagePath

	storageSubpath, exists := tusMeta["storage_subpath"]
	storageSubpath = strings.TrimSpace(storageSubpath)
	if !exists {
		return metadata{}, handler.NewError(
			"INVALID_METADATA",
			"metadata storage_subpath does not exist",
			http.StatusBadRequest,
		)
	}
	res.StorageSubpath = storageSubpath

	directoryPath, exists := tusMeta["directory_path"]
	directoryPath = strings.TrimSpace(directoryPath)
	if !exists {
		return metadata{}, handler.NewError(
			"INVALID_METADATA",
			"metadata directory_path does not exist",
			http.StatusBadRequest,
		)
	}
	res.DirectoryPath = directoryPath

	filename, exists := tusMeta["filename"]
	filename = strings.TrimSpace(filename)
	if !exists || filename == "" {
		return metadata{}, handler.NewError(
			"INVALID_METADATA",
			"metadata filename must not empty",
			http.StatusBadRequest,
		)
	}
	res.Filename = filename

	return res, nil
}

func (r *UploadRoute) Prefix() string        { return r.prefix }
func (r *UploadRoute) Handler() http.Handler { return r.handler }
