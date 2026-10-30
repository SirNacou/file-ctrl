package api

import (
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/labstack/echo/v5"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/app"
)

func Register(g huma.API, h *app.Handlers) {
	ListDirectory(g, h.ListDirectory)
	ListUnregisteredSources(g, h.ListUnregisteredSources)
	ListStorageSources(g, h.ListStorageSources)
	RegisterStorageSource(g, h.RegisterStorageSource)
	RemoveStorageSource(g, h.RemoveStorageSource)
	RenameStorage(g, h.RenameSource)
	CreateFolder(g, h.CreateFolder)
	RenameFile(g, h.RenameFile)
	DeleteFiles(g, h.DeleteFiles)
	DownloadFile(g, h.OpenFile)
	DownloadArchive(g, h.PrepareArchive)
}

// MountUpload keeps the Echo-specific prefix handling next to the
// upload handler instead of in the module.
func MountUpload(e *echo.Echo, endpoint *StorageUploadFilesEndpoint) {
	prefix := strings.TrimSuffix(endpoint.path, "/")
	wrapped := echo.WrapHandler(endpoint.handler)
	e.Any(prefix, wrapped)
	e.Any(prefix+"/*", wrapped)
}
