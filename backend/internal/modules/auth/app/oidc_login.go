package app

import (
	"context"
	"errors"
	"fmt"
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

type OIDCIdentity struct {
	Subject           string
	Email             string
	PreferredUsername string
}

var ErrOIDCExchange = errors.New("OIDC exchange and verification failed")

type OIDCIdentityVerifier interface {
	ExchangeAndVerify(ctx context.Context, code, expectedNonce string) (*OIDCIdentity, error)
}

type (
	OIDCLoginRequest struct {
		Code  string
		Nonce string
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
		verifier OIDCIdentityVerifier
	}
)

func NewOIDCLoginHandler(
	users OIDCLoginUsers,
	sessions OIDCLoginSessions,
	verifier OIDCIdentityVerifier,
) *OIDCLoginHandler {
	return &OIDCLoginHandler{
		users:    users,
		sessions: sessions,
		verifier: verifier,
	}
}

func (h *OIDCLoginHandler) Handle(
	ctx context.Context,
	req *OIDCLoginRequest,
) (*OIDCLoginResponse, error) {
	identity, err := h.verifier.ExchangeAndVerify(ctx, req.Code, req.Nonce)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrOIDCExchange, err)
	}

	user, err := h.users.FindByProvider(ctx, "oidc", identity.Subject)
	if err != nil {
		username := identity.PreferredUsername
		if username == "" {
			username = identity.Email
		}
		if username == "" {
			if len(identity.Subject) >= 8 {
				username = "user-" + identity.Subject[:8]
			} else {
				username = "user-" + identity.Subject
			}
		}

		role := domain.RoleUser
		userCount, _ := h.users.Count(ctx)
		if userCount == 0 {
			role = domain.RoleAdmin
		}

		user, err = domain.NewFederatedUser(username, identity.Subject, domain.ProviderOIDC, role)
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
