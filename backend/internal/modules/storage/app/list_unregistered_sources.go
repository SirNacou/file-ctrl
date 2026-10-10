package app

import (
	"context"
	"errors"
	"os"
	"slices"

	"github.com/sirnacou/file-ctrl/backend/internal/common/database/gen/models"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
	"github.com/stephenafamo/scan"
)

type ListUnregisteredSourcesRequest struct{}
type ListUnregisteredSourcesResponse struct {
	Dirs []string
}
type ListUnregisteredSourcesHandler struct {
	storageRoot string
	db          bob.Executor
}

func NewListUnregisteredSourcesHandler(storageRoot string) *ListUnregisteredSourcesHandler {
	return &ListUnregisteredSourcesHandler{storageRoot: storageRoot}
}

func (h *ListUnregisteredSourcesHandler) Handle(ctx context.Context, req *ListUnregisteredSourcesRequest) (*ListUnregisteredSourcesResponse, error) {
	entries, err := os.ReadDir(h.storageRoot)
	if err != nil {
		return nil, domain.NewErrFailToReadDir(h.storageRoot)
	}

	sourcePaths, err := h.getSourcePathsFromDB(ctx)
	if err != nil {
		return nil, err
	}

	dirs := make([]string, len(entries))
	for _, entry := range entries {
		if entry.IsDir() && !slices.Contains(sourcePaths, entry.Name()) {
			dirs = append(dirs, entry.Name())
		}
	}

	return &ListUnregisteredSourcesResponse{
		Dirs: dirs,
	}, nil
}

func (h *ListUnregisteredSourcesHandler) getSourcePathsFromDB(ctx context.Context) ([]string, error) {
	q := sqlite.Select(
		sm.Columns(models.StorageSources.Columns.Subpath.Name()),
		sm.From(models.StorageSources.Name()))

	sourcePaths, err := bob.All(ctx, h.db, q, scan.SingleColumnMapper[string])
	if err != nil {
		return nil, errors.New("failed to get source paths from database")
	}

	return sourcePaths, nil
}
