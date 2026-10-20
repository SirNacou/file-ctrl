package storage

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/common/localfs"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/adapters"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/ports/api"
	"github.com/stephenafamo/bob"
)

type StorageModule struct {
	listDirectoryHandler           *app.ListDirectoryHandler
	listUnregisteredSourcesHandler *app.ListUnregisteredSourcesHandler
	listStorageSourcesHandler      *app.ListStorageSourcesHandler
	registerStorageSourcesHandler  *app.RegisterStorageSourceHandler
	removeStorageSourceHandler     *app.RemoveStorageSourceHandler
	renameStorageHandler           *app.RenameStorageHandler
	createFolderHandler            *app.CreateFolderHandler
	renameFileHandler              *app.RenameFileHandler
	deleteFilesHandler             *app.DeleteFilesHandler
	downloadFileHandler            *app.DownloadFileHandler
}

func NewStorageModule(storageDir string, db bob.Executor) (*StorageModule, error) {
	storageRepo := adapters.NewStorageSourcesRepository(db)
	fsStore := localfs.NewStore(storageDir)

	return &StorageModule{
		listDirectoryHandler:           app.NewListDirectoryHandler(db, fsStore),
		listUnregisteredSourcesHandler: app.NewListUnregisteredSourcesHandler(db, fsStore),
		listStorageSourcesHandler:      app.NewListStorageSourcesHandler(db),
		registerStorageSourcesHandler: app.NewRegisterStorageSourceHandler(
			storageRepo,
			fsStore,
		),
		removeStorageSourceHandler: app.NewRemoveStorageSourceHandler(db),
		renameStorageHandler:       app.NewRenameStorageHandler(storageRepo, storageRepo),
		createFolderHandler:        app.NewCreateFolderHandler(storageRepo, fsStore),
		renameFileHandler:          app.NewRenameFileHandler(storageRepo, fsStore),
		deleteFilesHandler:         app.NewDeleteFilesHandler(storageRepo, fsStore),
		downloadFileHandler:        app.NewDownloadFileHandler(storageRepo, fsStore),
	}, nil
}

func (m *StorageModule) RegisterAPI(humaApi huma.API) {
	storageApi := huma.NewGroup(humaApi, "/storage")

	api.ListDirectory(storageApi, m.listDirectoryHandler)
	api.ListUnregisteredSources(storageApi, m.listUnregisteredSourcesHandler)
	api.ListStorageSources(storageApi, m.listStorageSourcesHandler)
	api.RegisterStorageSource(storageApi, m.registerStorageSourcesHandler)
	api.RemoveStorageSource(storageApi, m.removeStorageSourceHandler)
	api.RenameStorage(storageApi, m.renameStorageHandler)
	api.CreateFolder(storageApi, m.createFolderHandler)
	api.RenameFile(storageApi, m.renameFileHandler)
	api.DeleteFiles(storageApi, m.deleteFilesHandler)
	api.DownloadFile(storageApi, m.downloadFileHandler)
}
