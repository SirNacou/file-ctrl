package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/app"
)

func Register(api huma.API, h *app.Handlers, cookies *SessionCookies) {
	authAPI := huma.NewGroup(api, "/auth")

	GetStatus(authAPI, h.GetStatus)
	Login(authAPI, h.Login, cookies)
	Logout(authAPI, h.Logout, cookies)
	SetupAdmin(authAPI, h.SetupAdmin, cookies)
}
