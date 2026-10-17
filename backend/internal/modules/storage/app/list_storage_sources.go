package app

import (
	"context"
	"errors"
	"log/slog"
	"time"

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
	Subpath   string    `json:"subpath"`
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
		sources = append(sources, Source{
			Subpath:   s.Subpath,
			Name:      s.Name,
			CreatedAt: s.CreatedAt,
		})
	}

	return &ListStorageSourcesResponse{
		Sources: sources,
	}, nil
}
