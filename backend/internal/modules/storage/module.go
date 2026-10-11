package storage

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/ports/api"
	"github.com/stephenafamo/bob"
)

type StorageModule struct {
	listFilesHandler               *app.ListFilesHandler
	listUnregisteredSourcesHandler *app.ListUnregisteredSourcesHandler
	listStorageSourcesHandler      *app.ListStorageSourcesHandler
}

func NewStorageModule(storageDir string, db bob.Executor) (*StorageModule, error) {
	return &StorageModule{
		listFilesHandler:               app.NewListFilesHandler(storageDir),
		listUnregisteredSourcesHandler: app.NewListUnregisteredSourcesHandler(storageDir),
		listStorageSourcesHandler:      app.NewListStorageSourcesHandler(db),
	}, nil
}

func (m *StorageModule) RegisterAPI(humaApi huma.API) {
	api.ListFile(humaApi, m.listFilesHandler)
	api.ListUnregisteredSources(humaApi, m.listUnregisteredSourcesHandler)
	api.ListStorageSources(humaApi, m.listStorageSourcesHandler)
}
