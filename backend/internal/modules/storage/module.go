package storage

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type StorageModule struct {
}

func NewStorageModule() *StorageModule {
	return &StorageModule{}
}

func (m *StorageModule) RegisterAPI(api huma.API) {
	huma.Register(api, huma.Operation{
		Method:        "GET",
		Path:          "/list-files",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *struct{}) (*struct{}, error) {
		return nil, nil
	})
}
