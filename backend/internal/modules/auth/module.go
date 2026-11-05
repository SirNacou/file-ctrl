package auth

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/labstack/echo/v5"
	"github.com/sirnacou/file-ctrl/backend/internal/common/config"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/adapters"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/app"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/ports/api"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/ports/jobs"
	"github.com/stephenafamo/bob"
)

type AuthModule struct {
	handlers    *app.Handlers
	cookies     *api.SessionCookies
	oidcService *adapters.OIDCService
	oidcCookies *api.OIDCCookies
}

func NewAuthModule(ctx context.Context, db bob.Executor, env *config.Env) (*AuthModule, error) {
	userRepo := adapters.NewUserRepository(db)
	sessionRepo := adapters.NewSessionRepository(db)

	oidcService, err := adapters.NewOIDCService(ctx, env.OIDC)
	if err != nil {
		return nil, fmt.Errorf("create OIDC service: %w", err)
	}

	h := &app.Handlers{
		GetStatus:             app.NewGetStatusHandler(userRepo, sessionRepo, oidcService),
		Login:                 app.NewLoginHandler(userRepo, sessionRepo),
		Logout:                app.NewLogoutHandler(sessionRepo),
		OIDCLogin:             app.NewOIDCLoginHandler(userRepo, sessionRepo, oidcService),
		RemoveExpiredSessions: app.NewRemoveExpiredSessionsHandler(sessionRepo),
		SetupAdmin:            app.NewSetupAdminHandler(userRepo, sessionRepo),
		VerifySession:         app.NewVerifySessionHandler(userRepo, sessionRepo),
	}

	return &AuthModule{
		handlers:    h,
		cookies:     api.NewSessionCookies(env.CookieSecure()),
		oidcService: oidcService,
		oidcCookies: api.NewOIDCCookies(env.CookieSecure()),
	}, nil
}

func (m *AuthModule) RegisterAPI(humaAPI huma.API) {
	api.Register(humaAPI, m.handlers, m.cookies, m.oidcService, m.oidcCookies)
}

func (m *AuthModule) RegisterJobs(ctx context.Context) {
	jobs.CleanupSession(ctx, m.handlers.RemoveExpiredSessions)
}

func (m *AuthModule) Middleware() echo.MiddlewareFunc {
	return api.AuthMiddleware(m.handlers.VerifySession)
}
