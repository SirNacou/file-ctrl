package app

import (
	"context"
	"errors"

	"github.com/sirnacou/file-ctrl/backend/internal/common/database/gen/models"
	"github.com/stephenafamo/bob"
)

type RemoveStorageSourceHandler struct {
	db bob.Executor
}
type RemoveStorageSourceRequest struct {
	ID string
}
type RemoveStorageSourceResponse struct{}

func NewRemoveStorageSourceHandler(db bob.Executor) *RemoveStorageSourceHandler {
	return &RemoveStorageSourceHandler{db: db}
}

func (h *RemoveStorageSourceHandler) Handle(ctx context.Context, req *RemoveStorageSourceRequest) (*RemoveStorageSourceResponse, error) {
	_, err := models.StorageSources.Delete(models.DeleteWhere.StorageSources.ID.EQ(req.ID)).Exec(ctx, h.db)
	if err != nil {
		return nil, errors.New("failed to remove source")
	}

	return &RemoveStorageSourceResponse{}, nil
}
