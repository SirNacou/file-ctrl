package app

import (
	"context"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/domain"
)

type SessionDeleter interface {
	Delete(ctx context.Context, tokenHash string) error
}

type LogoutRequest struct {
	RawSessionToken string
}

type LogoutResponse struct{}

type LogoutHandler struct {
	sessions SessionDeleter
}

func (h *LogoutHandler) Handle(ctx context.Context, req *LogoutRequest) (*LogoutResponse, error) {
	if req.RawSessionToken == "" {
		return &LogoutResponse{}, nil
	}

	tokenHash := domain.HashSessionToken(req.RawSessionToken)
	if err := h.sessions.Delete(ctx, tokenHash); err != nil {
		return nil, err
	}

	return &LogoutResponse{}, nil
}
