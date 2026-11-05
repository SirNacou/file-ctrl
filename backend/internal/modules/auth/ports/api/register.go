package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/adapters"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/app"
)

func Register(
	api huma.API,
	h *app.Handlers,
	cookies *SessionCookies,
	oidcService *adapters.OIDCService,
	oidcCookies *OIDCCookies,
) {
	authAPI := huma.NewGroup(api, "/auth")
	authAPI.OpenAPI().Tags = []*huma.Tag{{Name: "Auth"}}

	GetStatus(authAPI, h.GetStatus)
	Login(authAPI, h.Login, cookies)
	Logout(authAPI, h.Logout, cookies)
	SetupAdmin(authAPI, h.SetupAdmin, cookies)
	OIDCLogin(authAPI, oidcService, oidcCookies)
	OIDCCallback(authAPI, h.OIDCLogin, oidcCookies, cookies)
}
