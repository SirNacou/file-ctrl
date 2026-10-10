package app

import (
	"context"
	"errors"
	"time"

	"github.com/sirnacou/file-ctrl/backend/internal/common/database/gen/queries"
	"github.com/stephenafamo/bob"
)

type ListStorageSourcesRequest struct{}
type ListStorageSourcesResponse struct {
	Sources []Source
}
type Source struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Subpath   string    `json:"sub_path"`
	CreatedAt time.Time `json:"created_at"`
}

type ListStorageSourcesHandler struct {
	db bob.Executor
}

func NewListStorageSourcesHandler(db bob.Executor) *ListStorageSourcesHandler {
	return &ListStorageSourcesHandler{db: db}
}

func (h *ListStorageSourcesHandler) Handle(ctx context.Context, req *ListStorageSourcesRequest) (*ListStorageSourcesResponse, error) {
	res, err := queries.AllStorageSources().All(ctx, h.db)
	if err != nil {
		return nil, errors.New("failed to query database")
	}

	sources := make([]Source, len(res))
	for _, s := range res {
		sources = append(sources, Source{
			ID:        s.ID,
			Name:      s.Name,
			Subpath:   s.Subpath,
			CreatedAt: time.Unix(s.CreatedAt.IntPart(), 0),
		})
	}

	return &ListStorageSourcesResponse{
		Sources: sources,
	}, nil
}
