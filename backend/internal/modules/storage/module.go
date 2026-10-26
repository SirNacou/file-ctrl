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
	storageDir string

	listDirectoryHandler           *app.ListDirectoryHandler
	listUnregisteredSourcesHandler *app.ListUnregisteredSourcesHandler
	listStorageSourcesHandler      *app.ListStorageSourcesHandler
	registerStorageSourcesHandler  *app.RegisterStorageSourceHandler
	removeStorageSourceHandler     *app.RemoveStorageSourceHandler
	renameStorageHandler           *app.RenameStorageHandler
	createFolderHandler            *app.CreateFolderHandler
	renameFileHandler              *app.RenameFileHandler
	deleteFilesHandler             *app.DeleteFilesHandler
	openFileHandler                *app.OpenFileHandler
	prepareArchiveHandler          *app.PrepareArchiveHandler
	validateUploadHandler          *app.ValidateUploadHandler
	uploadFileHandler              *app.UploadFileHandler
}

func NewStorageModule(storageDir string, db bob.Executor) (*StorageModule, error) {
	storageRepo := adapters.NewStorageSourcesRepository(db)
	localFileStorage := localfs.NewStorage(storageDir)

	return &StorageModule{
		storageDir: storageDir,

		listDirectoryHandler:           app.NewListDirectoryHandler(db, localFileStorage),
		listUnregisteredSourcesHandler: app.NewListUnregisteredSourcesHandler(db, localFileStorage),
		listStorageSourcesHandler:      app.NewListStorageSourcesHandler(db),
		registerStorageSourcesHandler: app.NewRegisterStorageSourceHandler(
			storageRepo,
			localFileStorage,
		),
		removeStorageSourceHandler: app.NewRemoveStorageSourceHandler(db),
		renameStorageHandler:       app.NewRenameStorageHandler(storageRepo, storageRepo),
		createFolderHandler:        app.NewCreateFolderHandler(storageRepo, localFileStorage),
		renameFileHandler:          app.NewRenameFileHandler(storageRepo, localFileStorage),
		deleteFilesHandler:         app.NewDeleteFilesHandler(storageRepo, localFileStorage),
		openFileHandler:            app.NewOpenFileHandler(storageRepo, localFileStorage),
		prepareArchiveHandler:      app.NewPrepareArchiveHandler(localFileStorage, storageRepo),
		validateUploadHandler:      app.NewValidateUploadHandler(storageRepo),
		uploadFileHandler:          app.NewUploadFileHandler(storageRepo, localFileStorage),
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
	api.DownloadFile(storageApi, m.openFileHandler)
	api.DownloadArchive(storageApi, m.prepareArchiveHandler)
}

func (m *StorageModule) UploadRoute() (*api.UploadRoute, error) {
	return api.NewUploadRoute(
		m.storageDir,
		m.validateUploadHandler,
		m.uploadFileHandler,
	)
}
