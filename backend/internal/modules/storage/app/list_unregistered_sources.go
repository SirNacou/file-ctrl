package app

import (
	"context"
	"errors"
	"log/slog"

	"github.com/sirnacou/file-ctrl/backend/internal/common/database/gen/models"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
	"github.com/stephenafamo/scan"
)

type (
	ListUnregisteredSourcesRequest  struct{}
	ListUnregisteredSourcesResponse struct {
		Dirs []string
	}
)

type ListUnregisteredSourcesHandler struct {
	db               bob.Executor
	storageDirLister FSDirLister
}

func NewListUnregisteredSourcesHandler(
	db bob.Executor,
	storageDirLister FSDirLister,
) *ListUnregisteredSourcesHandler {
	return &ListUnregisteredSourcesHandler{
		db:               db,
		storageDirLister: storageDirLister,
	}
}

func (h *ListUnregisteredSourcesHandler) Handle(
	ctx context.Context,
	req *ListUnregisteredSourcesRequest,
) (*ListUnregisteredSourcesResponse, error) {
	dirs, err := h.storageDirLister.ListDirs()
	if err != nil {
		slog.ErrorContext(ctx, "failed to read storage directory", slog.Any("error", err))
		return nil, domain.NewErrFailToReadDir("storage root")
	}

	registered, err := h.getSourcePathsFromDB(ctx)
	if err != nil {
		return nil, err
	}

	taken := make(map[string]struct{}, len(registered))
	for _, p := range registered {
		taken[p] = struct{}{}
	}

	unregistered := make([]string, 0, len(dirs))
	for _, d := range dirs {
		if _, ok := taken[d]; !ok {
			unregistered = append(unregistered, d)
		}
	}

	return &ListUnregisteredSourcesResponse{
		Dirs: dirs,
	}, nil
}

func (h *ListUnregisteredSourcesHandler) getSourcePathsFromDB(
	ctx context.Context,
) ([]string, error) {
	q := sqlite.Select(
		sm.Columns(models.StorageSources.Columns.Path.Name()),
		sm.From(models.StorageSources.Name()))

	sourcePaths, err := bob.All(ctx, h.db, q, scan.SingleColumnMapper[string])
	if err != nil {
		return nil, errors.New("failed to get source paths from database")
	}

	return sourcePaths, nil
}
