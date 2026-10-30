package storage

import (
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/labstack/echo/v5"
	"github.com/sirnacou/file-ctrl/backend/internal/common/localfs"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/adapters"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/ports/api"
	"github.com/stephenafamo/bob"
)

type StorageModule struct {
	handlers *app.Handlers
	upload   *api.StorageUploadFilesEndpoint
}

func NewStorageModule(storageDir string, db bob.Executor) (*StorageModule, error) {
	repo := adapters.NewStorageSourcesRepository(db)
	fs := localfs.NewStorage(storageDir)

	h := &app.Handlers{
		ListDirectory:           app.NewListDirectoryHandler(db, fs),
		ListUnregisteredSources: app.NewListUnregisteredSourcesHandler(db, fs),
		ListStorageSources:      app.NewListStorageSourcesHandler(db),
		RegisterStorageSource:   app.NewRegisterStorageSourceHandler(repo, fs),
		RemoveStorageSource:     app.NewRemoveStorageSourceHandler(db),
		RenameSource:            app.NewRenameStorageHandler(repo),
		CreateFolder:            app.NewCreateFolderHandler(repo, fs),
		RenameFile:              app.NewRenameFileHandler(repo, fs),
		DeleteFiles:             app.NewDeleteFilesHandler(repo, fs),
		OpenFile:                app.NewOpenFileHandler(repo, fs),
		PrepareArchive:          app.NewPrepareArchiveHandler(fs, repo),
		ValidateUpload:          app.NewValidateUploadHandler(repo),
		UploadFile:              app.NewUploadFileHandler(repo, fs),
	}

	upload, err := api.NewStorageUploadFilesEndpoint(storageDir, h.ValidateUpload, h.UploadFile)
	if err != nil {
		return nil, fmt.Errorf("create upload handler: %w", err)
	}

	return &StorageModule{handlers: h, upload: upload}, nil
}

// Can no longer fail: everything fallible happened in the constructor.
func (m *StorageModule) RegisterAPI(humaAPI huma.API, e *echo.Echo) {
	api.Register(huma.NewGroup(humaAPI, "/storage"), m.handlers)
	api.MountUpload(e, m.upload)
}
