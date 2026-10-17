package app

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/sirnacou/file-ctrl/backend/internal/common/database/gen/models"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
	"github.com/stephenafamo/bob"
)

type ListDirectoryRequest struct {
	Path           string
	StorageSubpath string
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
	db         bob.Executor
	storageDir string
}

func NewListDirectoryHandler(db bob.Executor, storageDir string) *ListDirectoryHandler {
	return &ListDirectoryHandler{
		db:         db,
		storageDir: storageDir,
	}
}

func (h *ListDirectoryHandler) Handle(
	ctx context.Context,
	req *ListDirectoryRequest,
) (*ListDirectoryResponse, error) {
	source, err := models.FindStorageSource(ctx, h.db, req.StorageSubpath)
	if err != nil {
		return nil, domain.ErrStorageSourceNotFound
	}

	subpath, err := domain.NewSourceSubpath(source.Subpath, false)
	if err != nil {
		return nil, domain.ErrStorageSourceNotFound
	}

	userPath, err := normalizeUserPath(req.Path)
	if err != nil {
		return nil, domain.ErrInvalidPath
	}

	dir, err := h.openDir(ctx, subpath.String(), userPath)
	if err != nil {
		return nil, domain.ErrInvalidPath
	}
	defer func() {
		err = dir.Close()
		if err != nil {
			slog.ErrorContext(ctx, "failed to close the dir", slog.Any("error", err))
		}
	}()

	items, err := readItems(dir, userPath)
	if err != nil {
		return nil, domain.ErrInvalidPath
	}

	return &ListDirectoryResponse{
		StorageName: source.Name,
		CurrentPath: displayPath(userPath),
		Items:       items,
	}, nil
}

// openDir chains storageDir -> source subpath -> user path.
// Each level is an os.Root, so nothing can escape the one above it.
func (h *ListDirectoryHandler) openDir(
	ctx context.Context,
	subpath, userPath string,
) (*os.File, error) {
	outer, err := os.OpenRoot(h.storageDir)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = outer.Close()
		if err != nil {
			slog.ErrorContext(ctx, "failed to open root dir", slog.Any("error", err))
		}
	}()

	src, err := outer.OpenRoot(filepath.FromSlash(subpath))
	if err != nil {
		return nil, err
	}
	defer func() {
		err = src.Close()
		if err != nil {
			slog.ErrorContext(ctx, "failed to open sub dir", slog.Any("error", err))
		}
	}()

	return src.Open(filepath.FromSlash(userPath))
}

func readItems(dir *os.File, userPath string) ([]DirectoryEntry, error) {
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

		items = append(items, toFileItem(userPath, e, info))
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

func toFileItem(parent string, e fs.DirEntry, info fs.FileInfo) DirectoryEntry {
	return DirectoryEntry{
		ModTime: info.ModTime(),
		Name:    e.Name(),
		Path:    path.Join(parent, e.Name()), // path.Join(".", "a") == "a"
		Size:    info.Size(),
		IsDir:   e.IsDir(),
	}
}

// normalizeUserPath validates the shape of the user path. os.Root still
// enforces the boundary; this just rejects obviously bad input early.
func normalizeUserPath(raw string) (string, error) {
	raw = strings.TrimPrefix(raw, "/")
	if raw == "" {
		return ".", nil
	}

	if strings.ContainsRune(raw, 0) || strings.Contains(raw, `\`) {
		return "", domain.ErrInvalidPath
	}

	if !filepath.IsLocal(filepath.FromSlash(raw)) {
		return "", domain.ErrInvalidPath
	}

	return path.Clean(raw), nil
}

func displayPath(userPath string) string {
	if userPath == "." {
		return ""
	}

	return userPath
}
