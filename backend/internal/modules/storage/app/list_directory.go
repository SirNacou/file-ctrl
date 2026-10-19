package app

import (
	"context"
	"io/fs"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/sirnacou/file-ctrl/backend/internal/common/database/gen/models"
	"github.com/sirnacou/file-ctrl/backend/internal/common/utils"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
	"github.com/stephenafamo/bob"
)

type ListDirectoryRequest struct {
	Subpath     string
	StoragePath string
}
type ListDirectoryResponse struct {
	StorageName string           `json:"storage_name"`
	CurrentPath string           `json:"current_path"`
	Items       []DirectoryEntry `json:"items"`
}

type DirectoryEntry struct {
	ModTime time.Time `json:"mod_time"`
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	Size    int64     `json:"size"`
	IsDir   bool      `json:"is_dir"`
}

type ListDirectoryHandler struct {
	db           bob.Executor
	sourceOpener SourceOpener
}

func NewListDirectoryHandler(db bob.Executor, sourceOpener SourceOpener) *ListDirectoryHandler {
	return &ListDirectoryHandler{
		db:           db,
		sourceOpener: sourceOpener,
	}
}

func (h *ListDirectoryHandler) Handle(
	ctx context.Context,
	req *ListDirectoryRequest,
) (*ListDirectoryResponse, error) {
	source, err := models.FindStorageSource(ctx, h.db, req.StoragePath)
	if err != nil {
		return nil, domain.ErrStorageSourceNotFound
	}

	entryPath, err := domain.NewEntryPath(req.Subpath)
	if err != nil {
		return nil, domain.ErrInvalidPath
	}

	root, err := h.sourceOpener.OpenSource(source.Path)
	if err != nil {
		return nil, domain.ErrStorageSourceNotFound
	}
	defer utils.CloseLog(ctx, root, "source root")

	dir, err := root.Open(entryPath.FS())
	if err != nil {
		return nil, domain.ErrInvalidPath
	}
	defer utils.CloseLog(ctx, dir, "directory")

	items, err := readItems(dir, entryPath)
	if err != nil {
		return nil, domain.ErrInvalidPath
	}

	return &ListDirectoryResponse{
		StorageName: source.Name,
		CurrentPath: entryPath.Display(),
		Items:       items,
	}, nil
}

func readItems(dir *os.File, entryPath domain.EntryPath) ([]DirectoryEntry, error) {
	entries, err := dir.ReadDir(-1) // fails if dir is not a directory
	if err != nil {
		return nil, err
	}

	items := make([]DirectoryEntry, 0, len(entries))
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue // removed between ReadDir and Info
		}

		items = append(items, toFileItem(entryPath, e, info))
	}

	// Directories first, then by name
	slices.SortFunc(items, func(a, b DirectoryEntry) int {
		if a.IsDir != b.IsDir {
			if a.IsDir {
				return -1
			}

			return 1
		}

		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})

	return items, nil
}

func toFileItem(parent domain.EntryPath, e fs.DirEntry, info fs.FileInfo) DirectoryEntry {
	return DirectoryEntry{
		ModTime: info.ModTime(),
		Name:    e.Name(),
		Path:    parent.Child(e.Name()),
		Size:    info.Size(),
		IsDir:   e.IsDir(),
	}
}
