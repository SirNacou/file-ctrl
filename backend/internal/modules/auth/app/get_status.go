package app

import (
	"context"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/domain"
	"uuid"
)

type (
	UserFinderCounter interface {
		Count(ctx context.Context) (int64, error)
		Find(ctx context.Context, ID uuid.UUID) (*domain.User, error)
	}

	// OIDCStatusChecker determines if SSO is configured on the backend
	OIDCStatusChecker interface {
		IsEnabled() bool
	}

	GetStatusRequest struct {
		RawSessionToken string
	}

	GetStatusResponse struct {
		NeedsSetup    bool
		Authenticated bool
		OidcEnabled   bool // <--- Added
		Username      string
		Role          domain.Role
	}

	GetStatusHandler struct {
		userFinderCounter UserFinderCounter
		sessionFinder     SessionFinder
		oidcChecker       OIDCStatusChecker // <--- Injected
	}
)

func NewGetStatusHandler(
	userFinderCounter UserFinderCounter,
	sessionFinder SessionFinder,
	oidcChecker OIDCStatusChecker, // <--- Pass your OIDCService or any type implementing IsEnabled() bool
) *GetStatusHandler {
	return &GetStatusHandler{
		userFinderCounter: userFinderCounter,
		sessionFinder:     sessionFinder,
		oidcChecker:       oidcChecker,
	}
}

func (h *GetStatusHandler) Handle(
	ctx context.Context,
	req *GetStatusRequest,
) (*GetStatusResponse, error) {
	oidcEnabled := h.oidcChecker != nil && h.oidcChecker.IsEnabled()

	count, err := h.userFinderCounter.Count(ctx)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return &GetStatusResponse{
			NeedsSetup:    true,
			Authenticated: false,
			OidcEnabled:   oidcEnabled,
		}, nil
	}

	if req.RawSessionToken == "" {
		return &GetStatusResponse{
			NeedsSetup:    false,
			Authenticated: false,
			OidcEnabled:   oidcEnabled,
		}, nil
	}

	tokenHash := domain.HashSessionToken(req.RawSessionToken)
	session, err := h.sessionFinder.FindByTokenHash(ctx, tokenHash)
	if err != nil || session == nil || session.IsExpired() {
		return &GetStatusResponse{
			NeedsSetup:    false,
			Authenticated: false,
			OidcEnabled:   oidcEnabled,
		}, nil
	}

	user, err := h.userFinderCounter.Find(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	return &GetStatusResponse{
		NeedsSetup:    false,
		Authenticated: true,
		OidcEnabled:   oidcEnabled,
		Username:      user.Username,
		Role:          user.Role,
	}, nil
}
