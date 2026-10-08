package app

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type ListFilesRequest struct {
	Path string
}
type ListFilesResponse struct {
	CurrentPath string     `json:"current_path"`
	Items       []FileItem `json:"items"`
}

type FileItem struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	IsDir   bool      `json:"is_dir"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
}

type ListFilesHandler struct {
	storageRoot string
}

func NewListFilesHandler(storageRoot string) *ListFilesHandler {
	return &ListFilesHandler{
		storageRoot: storageRoot,
	}
}

func (h *ListFilesHandler) Handle(ctx context.Context, req *ListFilesRequest) (*ListFilesResponse, error) {
	if req.Path == "" {
		req.Path = "/"
	}
	if !strings.HasPrefix(req.Path, "/") {
		req.Path = "/" + req.Path
	}

	diskPath := filepath.Join(h.storageRoot, filepath.FromSlash(req.Path))

	entries, err := os.ReadDir(diskPath)
	if err != nil {
		return nil, domain.NewErrFailToReadDir(diskPath)
	}

	items := make([]FileItem, 0, len(entries))
	for _, v := range entries {
		info, err := v.Info()
		if err != nil {
			continue
		}

		var size int64 = 0
		if !v.IsDir() {
			size = info.Size()
		}

		// Use req.Path so the frontend receives virtual URLs, not host disk paths
		itemPath := path.Join(req.Path, v.Name())
		items = append(items, FileItem{
			Name:    v.Name(),
			Path:    itemPath,
			IsDir:   v.IsDir(),
			Size:    size,
			ModTime: info.ModTime().UTC(),
		})
	}

	return &ListFilesResponse{
		CurrentPath: req.Path,
		Items:       items,
	}, nil
}
