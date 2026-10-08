package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/common/config"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/ports/api"
)

type StorageModule struct {
	listFilesHandler *app.ListFilesHandler
}

func NewStorageModule(cfg *config.Config) (*StorageModule, error) {
	rootDir, err := createStorageDirectory(cfg.StorageRoot())
	if err != nil {
		return nil, err
	}

	return &StorageModule{
		listFilesHandler: app.NewListFilesHandler(rootDir),
	}, nil
}

func (m *StorageModule) RegisterAPI(humaApi huma.API) {
	api.ListFile(humaApi, m.listFilesHandler)
}

func createStorageDirectory(rootPath string) (string, error) {
	absRoot, err := filepath.Abs(filepath.Clean(rootPath))
	if err != nil {
		return "", fmt.Errorf("invalid storage root path: %w", err)
	}

	if err := os.MkdirAll(absRoot, 0755); err != nil {
		return "", fmt.Errorf("failed to create storage directory %q: %w", absRoot, err)
	}

	return absRoot, nil
}
