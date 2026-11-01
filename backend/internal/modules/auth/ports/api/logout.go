package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/app"
)

type LogoutReponse struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
}

func Logout(api huma.API, h *app.LogoutHandler, cookies *SessionCookies) {
	huma.Register(api, huma.Operation{
		OperationID: "auth-logout",
		Method:      http.MethodPost,
		Path:        "/logout",
		Summary:     "Invalidate session and clear cookie",
	}, func(ctx context.Context, input *struct{}) (*LogoutReponse, error) {
		hCtx := ctx.(huma.Context)
		cookie, _ := huma.ReadCookie(hCtx, SessionCookieName)
		rawToken := ""
		if cookie != nil {
			rawToken = cookie.Value
		}

		_, _ = h.Handle(ctx, &app.LogoutRequest{RawSessionToken: rawToken})

		clearCookie := cookies.Clear()

		return &LogoutReponse{
			SetCookie: *clearCookie,
		}, nil
	})
}
