package app

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/sirnacou/file-ctrl/backend/internal/common/database"
	"github.com/sirnacou/file-ctrl/backend/internal/common/database/gen/models"
	"github.com/stephenafamo/bob"
)

type (
	ListStorageSourcesRequest  struct{}
	ListStorageSourcesResponse struct {
		Sources []Source
	}
)

type Source struct {
	CreatedAt time.Time `json:"created_at"`
	Path      string    `json:"path"`
	Name      string    `json:"name"`
}

type ListStorageSourcesHandler struct {
	db bob.Executor
}

func NewListStorageSourcesHandler(db bob.Executor) *ListStorageSourcesHandler {
	return &ListStorageSourcesHandler{db: db}
}

func (h *ListStorageSourcesHandler) Handle(
	ctx context.Context,
	req *ListStorageSourcesRequest,
) (*ListStorageSourcesResponse, error) {
	res, err := models.StorageSources.Query().All(ctx, h.db)
	if err != nil {
		slog.ErrorContext(ctx, "failed to query storage sources", slog.String("error", err.Error()))
		return nil, errors.New("failed to query database")
	}

	sources := make([]Source, 0, len(res))
	for _, s := range res {
		createdAt, err := database.ParseTime(s.CreatedAt)
		if err != nil {
			slog.ErrorContext(
				ctx,
				"failed to parse storage source creation time",
				slog.String("error", err.Error()),
			)
			return nil, errors.New("failed to parse database result")
		}

		sources = append(sources, Source{
			Path:      s.Path,
			Name:      s.Name,
			CreatedAt: createdAt,
		})
	}

	return &ListStorageSourcesResponse{
		Sources: sources,
	}, nil
}
