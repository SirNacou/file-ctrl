package app

import (
	"context"
	"time"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/domain"
	"uuid"
)

type OIDCLoginUsers interface {
	Count(ctx context.Context) (int64, error)
	FindByProvider(ctx context.Context, provider string, providerID string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
}

type OIDCLoginSessions interface {
	Create(ctx context.Context, session *domain.Session) error
}

type (
	OIDCLoginRequest struct {
		Subject           string
		Email             string
		PreferredUsername string
	}
	OIDCLoginResponse struct {
		ID        uuid.UUID
		Username  string
		Role      domain.Role
		RawToken  string
		ExpiresAt time.Time
	}
	OIDCLoginHandler struct {
		users    OIDCLoginUsers
		sessions OIDCLoginSessions
	}
)

func NewOIDCLoginHandler(
	users OIDCLoginUsers,
	sessions OIDCLoginSessions,
) *OIDCLoginHandler {
	return &OIDCLoginHandler{
		users:    users,
		sessions: sessions,
	}
}

func (h *OIDCLoginHandler) Handle(
	ctx context.Context,
	req *OIDCLoginRequest,
) (*OIDCLoginResponse, error) {
	user, err := h.users.FindByProvider(ctx, "oidc", req.Subject)
	if err != nil {
		username := req.PreferredUsername
		if username == "" {
			username = req.Email
		}
		if username == "" {
			username = "user-" + req.Subject[:8]
		}

		role := domain.RoleUser
		userCount, _ := h.users.Count(ctx)
		if userCount == 0 {
			role = domain.RoleAdmin
		}

		user, err = domain.NewFederatedUser(username, req.Subject, domain.ProviderOIDC, role)
		if err != nil {
			return nil, err
		}

		if err := h.users.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	sess, token, err := domain.NewSession(user.ID)
	if err != nil {
		return nil, err
	}

	if err := h.sessions.Create(ctx, sess); err != nil {
		return nil, err
	}

	return &OIDCLoginResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		RawToken:  token.Raw,
		ExpiresAt: sess.ExpiresAt,
	}, nil
}
