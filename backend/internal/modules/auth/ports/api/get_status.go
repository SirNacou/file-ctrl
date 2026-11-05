package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/app"
)

type AuthStatusRequest struct {
	RawToken string `cookie:"filectrl_session"`
}

type AuthStatusResponse struct {
	Body AuthStatusResponseBody
}

type AuthStatusResponseBody struct {
	NeedsSetup    bool   `json:"needs_setup"`
	Authenticated bool   `json:"authenticated"`
	OidcEnabled   bool   `json:"oidc_enabled"` // <--- Added
	Username      string `json:"username,omitempty"`
	Role          string `json:"role,omitempty"`
}

func GetStatus(api huma.API, h *app.GetStatusHandler) {
	huma.Register(api, huma.Operation{
		OperationID: "auth-get-status",
		Method:      http.MethodGet,
		Path:        "/status",
		Summary:     "Get current authentication state and setup status",
	}, func(ctx context.Context, input *AuthStatusRequest) (*AuthStatusResponse, error) {
		res, err := h.Handle(
			ctx,
			&app.GetStatusRequest{RawSessionToken: input.RawToken},
		)
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to check status", err)
		}

		return &AuthStatusResponse{
			Body: AuthStatusResponseBody{
				NeedsSetup:    res.NeedsSetup,
				Authenticated: res.Authenticated,
				OidcEnabled:   res.OidcEnabled, // <--- Mapped
				Username:      res.Username,
				Role:          string(res.Role),
			},
		}, nil
	})
}
