package app

import (
	"context"
	"uuid"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/domain"
)

type (
	UserFinderCounter interface {
		Count(ctx context.Context) (int64, error)
		Find(ctx context.Context, ID uuid.UUID) (*domain.User, error)
	}
	GetStatusRequest struct {
		RawSessionToken string
	}
	GetStatusResponse struct {
		NeedsSetup    bool
		Authenticated bool
		User          *domain.User
	}
	GetStatusHandler struct {
		sessionFinder     SessionFinder
		userFinderCounter UserFinderCounter
	}
)

func NewGetStatusHandler(
	sessionFinder SessionFinder,
	userFinderCounter UserFinderCounter,
) *GetStatusHandler {
	return &GetStatusHandler{
		sessionFinder:     sessionFinder,
		userFinderCounter: userFinderCounter,
	}
}

func (h *GetStatusHandler) Handle(
	ctx context.Context,
	req *GetStatusRequest,
) (*GetStatusResponse, error) {
	count, err := h.userFinderCounter.Count(ctx)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return &GetStatusResponse{
			NeedsSetup:    true,
			Authenticated: false,
		}, nil
	}

	if req.RawSessionToken == "" {
		return &GetStatusResponse{
			NeedsSetup:    false,
			Authenticated: false,
		}, nil
	}

	tokenHash := domain.HashSessionToken(req.RawSessionToken)
	session, err := h.sessionFinder.FindByTokenHash(ctx, tokenHash)
	if err != nil || session.IsExpired() {
		return &GetStatusResponse{
			NeedsSetup:    false,
			Authenticated: false,
		}, nil
	}

	user, err := h.userFinderCounter.Find(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	return &GetStatusResponse{
		NeedsSetup:    false,
		Authenticated: true,
		User:          user,
	}, nil
}
