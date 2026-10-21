package app

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/sirnacou/file-ctrl/backend/internal/common/database/gen/models"
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
	db          bob.Executor
	fsDirReader DirReader
}

func NewListDirectoryHandler(db bob.Executor, fsDirReader DirReader) *ListDirectoryHandler {
	return &ListDirectoryHandler{
		db:          db,
		fsDirReader: fsDirReader,
	}
}

func (h *ListDirectoryHandler) Handle(
	ctx context.Context,
	req *ListDirectoryRequest,
) (*ListDirectoryResponse, error) {
	storageSource, err := models.FindStorageSource(ctx, h.db, req.StoragePath)
	if err != nil {
		return nil, domain.ErrStorageSourceNotFound
	}

	sourcePath, err := domain.NewSourcePathWithDefault(storageSource.Path)
	if err != nil {
		return nil, err
	}

	entryPath, err := domain.NewEntryPath(req.Subpath)
	if err != nil {
		return nil, domain.ErrInvalidPath
	}

	// Note: omit ctx here if FSDirReader.ReadDir matches localfs.Store.ReadDir(sourcePath, dir)
	entries, err := h.fsDirReader.ReadDir(ctx, sourcePath, entryPath)
	if err != nil {
		return nil, domain.ErrInvalidPath
	}

	items := readItems(entries, entryPath)

	return &ListDirectoryResponse{
		StorageName: storageSource.Name,
		CurrentPath: entryPath.Display(),
		Items:       items,
	}, nil
}

func readItems(entries []domain.FileEntry, entryPath domain.EntryPath) []DirectoryEntry {
	items := make([]DirectoryEntry, 0, len(entries))
	for _, e := range entries {
		items = append(items, toFileItem(entryPath, e))
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

	return items
}

func toFileItem(parent domain.EntryPath, e domain.FileEntry) DirectoryEntry {
	return DirectoryEntry{
		ModTime: e.ModTime,
		Name:    e.Name,
		Path:    parent.Child(e.Name),
		Size:    e.Size,
		IsDir:   e.IsDir,
	}
}
