package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/adapters"
)

type (
	OIDCLoginRequest  struct{}
	OIDCLoginResponse struct {
		Status    int           `status:"302"`
		Location  string        `header:"Location"`
		SetCookie []http.Cookie `header:"Set-Cookie"`
	}
)

func OIDCLogin(api huma.API, oidcService *adapters.OIDCService, cookies *OIDCCookies) {
	huma.Register(api, huma.Operation{
		OperationID: "auth-oidc-login",
		Method:      http.MethodGet,
		Path:        "/oidc/login",
		Summary:     "Initiate OIDC SSO redirect",
		Tags:        []string{"Auth"},
	}, func(ctx context.Context, i *OIDCLoginRequest) (*OIDCLoginResponse, error) {
		if !oidcService.IsEnabled() {
			return nil, huma.Error404NotFound("OIDC authentication is not configured")
		}

		state, nonce, err := oidcService.GenerateStateAndNonce()
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to generate security state", err)
		}

		authURL := oidcService.AuthCodeURL(state, nonce)

		stateCookie := cookies.NewState(state)

		nonceCookie := cookies.NewNonce(nonce)

		return &OIDCLoginResponse{
			Status:    http.StatusFound,
			Location:  authURL,
			SetCookie: []http.Cookie{stateCookie, nonceCookie},
		}, nil
	})
}
