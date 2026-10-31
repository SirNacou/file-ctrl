package app

import (
	"context"
	"uuid"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/domain"
)

type UserFinder interface {
	Find(
		ctx context.Context,
		ID uuid.UUID,
	) (*domain.User, error)
}

type (
	VerifySessionRequest struct {
		SessionToken string
	}
	VerifySessionResponse struct {
		UserID   uuid.UUID
		Username string
		Role     domain.Role
	}
	VerifySessionHandler struct {
		users    UserFinder
		sessions SessionFinder
	}
)

func NewVerifySessionHandler(users UserFinder, sessions SessionFinder) *VerifySessionHandler {
	return &VerifySessionHandler{
		users:    users,
		sessions: sessions,
	}
}

func (h *VerifySessionHandler) Handle(
	ctx context.Context,
	req *VerifySessionRequest,
) (*VerifySessionResponse, error) {
	token := domain.HashSessionToken(req.SessionToken)
	sess, err := h.sessions.FindByTokenHash(ctx, token)
	if err != nil {
		return nil, domain.ErrSessionNotFound
	}

	user, err := h.users.Find(ctx, sess.UserID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return &VerifySessionResponse{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	}, nil
}
