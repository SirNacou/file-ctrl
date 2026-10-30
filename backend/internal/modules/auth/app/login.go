package app

import (
	"context"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/domain"
)

type UserFinder interface {
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
}

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	User  *domain.User
	Token domain.SessionToken
}

type LoginHandler struct {
	users    UserFinder
	sessions SessionCreator
}

func NewLoginHandler(users UserFinder, sessions SessionCreator) *LoginHandler {
	return &LoginHandler{
		users:    users,
		sessions: sessions,
	}
}

func (h *LoginHandler) Handle(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	user, err := h.users.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	password, err := domain.NewPassword(req.Password)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	valid, err := password.VerifyPassword(user.PasswordHash)
	if err != nil || !valid {
		return nil, domain.ErrInvalidCredentials
	}

	sess, token, err := domain.NewSession(user.ID)
	if err != nil {
		return nil, err
	}

	if err := h.sessions.Create(ctx, sess); err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}
