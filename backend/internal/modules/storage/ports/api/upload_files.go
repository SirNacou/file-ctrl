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

type StorageUploadFilesEndpoint struct {
	handler http.Handler
	path    string
}

func NewStorageUploadFilesEndpoint(
	storageDir string,
	validateUploadHandler *app.ValidateUploadHandler,
	uploadFileHandler *app.UploadFileHandler,
) (*StorageUploadFilesEndpoint, error) {
	stagingDir := filepath.Join(storageDir, ".uploads")
	if err := os.MkdirAll(stagingDir, 0o755); err != nil {
		return nil, err
	}

	store := filestore.New(stagingDir)
	composer := handler.NewStoreComposer()
	store.UseIn(composer)

	tusHandler, err := handler.NewHandler(handler.Config{
		StoreComposer:             composer,
		BasePath:                  UploadBasePath,
		Logger:                    slog.Default(),
		PreUploadCreateCallback:   preUploadCreateHandler(validateUploadHandler),
		PreFinishResponseCallback: preFinishResponseHandler(uploadFileHandler),
	})
	if err != nil {
		return nil, err
	}

	return &StorageUploadFilesEndpoint{
		handler: tusHandler,
		path:    UploadBasePath,
	}, nil
}

func preUploadCreateHandler(
	validatehandler *app.ValidateUploadHandler,
) func(hook handler.HookEvent) (handler.HTTPResponse, handler.FileInfoChanges, error) {
	return func(hook handler.HookEvent) (handler.HTTPResponse, handler.FileInfoChanges, error) {
		slog.InfoContext(hook.Context, "Metadata", slog.Any("meta", hook.Upload.MetaData))
		meta, err := mapToMetadata(hook.Upload.MetaData)
		if err != nil {
			return handler.HTTPResponse{
				StatusCode: http.StatusBadRequest,
			}, handler.FileInfoChanges{}, err
		}

		_, err = validatehandler.Handle(hook.Context, &app.ValidateUploadRequest{
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
	}
}

func preFinishResponseHandler(
	uploadFileHandler *app.UploadFileHandler,
) func(hook handler.HookEvent) (handler.HTTPResponse, error) {
	return func(hook handler.HookEvent) (handler.HTTPResponse, error) {
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
	}
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
		return metadata{}, errInvalidMetadata("metadata storage_path must not empty")
	}
	res.StoragePath = storagePath

	storageSubpath, exists := tusMeta["storage_subpath"]
	storageSubpath = strings.TrimSpace(storageSubpath)
	if !exists {
		return metadata{}, errInvalidMetadata("metadata storage_subpath does not exist")
	}
	res.StorageSubpath = storageSubpath

	directoryPath, exists := tusMeta["directory_path"]
	directoryPath = strings.TrimSpace(directoryPath)
	if !exists {
		return metadata{}, errInvalidMetadata("metadata directory_path does not exist")
	}
	res.DirectoryPath = directoryPath

	filename, exists := tusMeta["filename"]
	filename = strings.TrimSpace(filename)
	if !exists || filename == "" {
		return metadata{}, errInvalidMetadata("metadata filename must not empty")
	}
	res.Filename = filename

	return res, nil
}

func errInvalidMetadata(message string) error {
	return handler.NewError(
		"INVALID_METADATA",
		message,
		http.StatusBadRequest,
	)
}
