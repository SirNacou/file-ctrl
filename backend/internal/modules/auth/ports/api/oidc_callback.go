package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/app"
)

type OIDCCallbackRequest struct {
	Code  string `query:"code"  required:"true" doc:"Authorization code from provider"`
	State string `query:"state" required:"true" doc:"OIDC state returned from provider"`

	CookieState string `cookie:"oidc_state" doc:"Expected state stored in client cookie"`
	CookieNonce string `cookie:"oidc_nonce" doc:"Nonce used during authorization request"`
}
type OIDCCallbackResponse struct {
	Status int `status:"302"`

	Location  string        `header:"Location"`
	SetCookie []http.Cookie `header:"Set-Cookie"`
}

func OIDCCallback(
	api huma.API,
	h *app.OIDCLoginHandler,
	oidcCookies *OIDCCookies,
	cookies *SessionCookies,
) {
	huma.Register(api, huma.Operation{
		OperationID: "auth-oidc-callback",
		Method:      http.MethodGet,
		Path:        "/oidc/callback",
		Summary:     "Exchange authorization code and establish session",
		Tags:        []string{"Auth"},
	}, func(ctx context.Context, i *OIDCCallbackRequest) (*OIDCCallbackResponse, error) {
		// 1. Verify CSRF state parameter matches the browser's cookie
		if i.CookieState == "" || i.CookieState != i.State {
			return nil, huma.Error400BadRequest("Invalid or missing state parameter")
		}

		if i.CookieNonce == "" {
			return nil, huma.Error400BadRequest("Missing nonce verification cookie")
		}

		// 2. Exchange/verify the OIDC code and create the local session
		res, err := h.Handle(ctx, &app.OIDCLoginRequest{
			Code:  i.Code,
			Nonce: i.CookieNonce,
		})
		if err != nil {
			if errors.Is(err, app.ErrOIDCExchange) {
				return nil, huma.Error401Unauthorized(err.Error())
			}
			return nil, huma.Error500InternalServerError("Failed to create user session", err)
		}

		// 4. Invalidate temporary state/nonce cookies
		clearStateCookie := oidcCookies.ClearState()

		clearNonceCookie := oidcCookies.ClearNonce()

		// 5. Build standard app session cookie (NewSessionCookie from Phase 1)
		sessionCookie := cookies.New(res.RawToken, res.ExpiresAt)

		return &OIDCCallbackResponse{
			Status:   http.StatusFound,
			Location: "/",
			SetCookie: []http.Cookie{
				clearStateCookie,
				clearNonceCookie,
				*sessionCookie,
			},
		}, nil
	})
}
