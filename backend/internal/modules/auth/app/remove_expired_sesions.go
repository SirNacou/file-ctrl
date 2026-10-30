package app

import (
	"context"
	"time"
)

type ExpiredSessionsDeleter interface {
	DeleteAllBefore(ctx context.Context, cutoff time.Time) error
}

type (
	RemoveExpiredSessionsRequest struct {
		Cutoff time.Time
	}
	RemoveExpiredSessionsResponse struct{}
	RemoveExpiredSessionsHandler  struct {
		sessions ExpiredSessionsDeleter
	}
)

func NewRemoveExpiredSessionsHandler(
	sessions ExpiredSessionsDeleter,
) *RemoveExpiredSessionsHandler {
	return &RemoveExpiredSessionsHandler{
		sessions: sessions,
	}
}

func (h *RemoveExpiredSessionsHandler) Handle(
	ctx context.Context,
	req *RemoveExpiredSessionsRequest,
) (*RemoveExpiredSessionsResponse, error) {
	err := h.sessions.DeleteAllBefore(ctx, req.Cutoff)
	if err != nil {
		return nil, err
	}

	return &RemoveExpiredSessionsResponse{}, nil
}
