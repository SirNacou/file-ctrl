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
		return nil, domain.NewErrInvalidPath(req.Subpath)
	}

	// Note: omit ctx here if FSDirReader.ReadDir matches localfs.Store.ReadDir(sourcePath, dir)
	entries, err := h.fsDirReader.ReadDir(ctx, sourcePath, entryPath)
	if err != nil {
		return nil, domain.NewErrInvalidPath(entryPath.String())
	}

	items, err := readItems(entries, entryPath)
	if err != nil {
		return nil, err
	}

	return &ListDirectoryResponse{
		StorageName: storageSource.Name,
		CurrentPath: entryPath.Display(),
		Items:       items,
	}, nil
}

func readItems(entries []domain.FileMeta, entryPath domain.EntryPath) ([]DirectoryEntry, error) {
	items := make([]DirectoryEntry, 0, len(entries))
	for _, e := range entries {
		item, err := toFileItem(entryPath, e)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
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

func toFileItem(parent domain.EntryPath, e domain.FileMeta) (DirectoryEntry, error) {
	path, err := parent.Join(e.Name)
	if err != nil {
		return DirectoryEntry{}, err
	}
	return DirectoryEntry{
		ModTime: e.ModTime,
		Name:    e.Name,
		Path:    path.Display(),
		Size:    e.Size,
		IsDir:   e.IsDir,
	}, nil
}
