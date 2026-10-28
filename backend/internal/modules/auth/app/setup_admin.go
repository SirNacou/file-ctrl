package app

import (
	"context"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/domain"
)

type UserCounterCreator interface {
	Create(
		ctx context.Context,
		user *domain.User,
	) error
	Count(ctx context.Context) (int64, error)
}

type SetupAdminRequest struct {
	Username string
	Password string
}

type SetupAdminResponse struct {
	User  *domain.User
	Token domain.SessionToken
}

type SetupAdminHandler struct {
	userRepo    UserCounterCreator
	sessionRepo SessionCreator
}

func NewSetupAdminHandler(
	userRepo UserCounterCreator,
	sessionRepo SessionCreator,
) *SetupAdminHandler {
	return &SetupAdminHandler{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (h *SetupAdminHandler) Handle(
	ctx context.Context,
	req *SetupAdminRequest,
) (*SetupAdminResponse, error) {
	count, err := h.userRepo.Count(ctx)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, domain.ErrSetupAlreadyDone
	}

	password, err := domain.NewPassword(req.Password)
	if err != nil {
		return nil, err
	}

	hash, err := password.HashPassword()
	if err != nil {
		return nil, err
	}

	user, err := domain.NewUser(req.Username, hash, domain.RoleAdmin)
	if err != nil {
		return nil, err
	}

	session, token, err := domain.NewSession(user.ID)
	if err != nil {
		return nil, err
	}

	if err := h.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return &SetupAdminResponse{
		User:  user,
		Token: token,
	}, nil
}
