package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/app"
)

type LogoutRequest struct {
	HumaSessionCookie
}
type LogoutResponse struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
}

func Logout(api huma.API, h *app.LogoutHandler, cookies *SessionCookies) {
	huma.Register(api, huma.Operation{
		OperationID: "auth-logout",
		Method:      http.MethodPost,
		Path:        "/logout",
		Summary:     "Invalidate session and clear cookie",
	}, func(ctx context.Context, input *LogoutRequest) (*LogoutResponse, error) {
		_, _ = h.Handle(ctx, &app.LogoutRequest{RawSessionToken: input.SessionCookie.Value})

		clearCookie := cookies.Clear()

		return &LogoutResponse{
			SetCookie: *clearCookie,
		}, nil
	})
}
