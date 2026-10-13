package storage

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/adapters"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/ports/api"
	"github.com/stephenafamo/bob"
)

type StorageModule struct {
	listFilesHandler               *app.ListDirectoryHandler
	listUnregisteredSourcesHandler *app.ListUnregisteredSourcesHandler
	listStorageSourcesHandler      *app.ListStorageSourcesHandler
	registerStorageSourcesHandler  *app.RegisterStorageSourceHandler
	removeStorageSourceHandler     *app.RemoveStorageSourceHandler
}

func NewStorageModule(storageDir string, db bob.Executor) (*StorageModule, error) {
	storageRepo := adapters.NewStorageSourcesRepository(db)

	return &StorageModule{
		listFilesHandler:               app.NewListDirectoryHandler(db, storageDir),
		listUnregisteredSourcesHandler: app.NewListUnregisteredSourcesHandler(db, storageDir),
		listStorageSourcesHandler:      app.NewListStorageSourcesHandler(db),
		registerStorageSourcesHandler: app.NewRegisterStorageSourceHandler(
			storageRepo,
			storageDir,
		),
		removeStorageSourceHandler: app.NewRemoveStorageSourceHandler(db),
	}, nil
}

func (m *StorageModule) RegisterAPI(humaApi huma.API) {
	storageApi := huma.NewGroup(humaApi, "/storage")

	api.ListDirectory(storageApi, m.listFilesHandler)
	api.ListUnregisteredSources(storageApi, m.listUnregisteredSourcesHandler)
	api.ListStorageSources(storageApi, m.listStorageSourcesHandler)
	api.RegisterStorageSource(storageApi, m.registerStorageSourcesHandler)
	api.RemoveStorageSource(storageApi, m.removeStorageSourceHandler)
}
