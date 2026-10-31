package auth

import (
	"context"

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
	handlers *app.Handlers
	cookies  *api.SessionCookies
}

func NewAuthModule(db bob.Executor, env *config.Env) *AuthModule {
	userRepo := adapters.NewUserRepository(db)
	sessionRepo := adapters.NewSessionRepository(db)

	h := &app.Handlers{
		GetStatus:             app.NewGetStatusHandler(userRepo, sessionRepo),
		Login:                 app.NewLoginHandler(userRepo, sessionRepo),
		Logout:                app.NewLogoutHandler(sessionRepo),
		RemoveExpiredSessions: app.NewRemoveExpiredSessionsHandler(sessionRepo),
		SetupAdmin:            app.NewSetupAdminHandler(userRepo, sessionRepo),
		VerifySession:         app.NewVerifySessionHandler(userRepo, sessionRepo),
	}

	return &AuthModule{
		handlers: h,
		cookies:  api.NewSessionCookies(env.CookieSecure()),
	}
}

func (m *AuthModule) RegisterAPI(humaAPI huma.API) {
	api.Register(humaAPI, m.handlers, m.cookies)
}

func (m *AuthModule) RegisterJobs(ctx context.Context) {
	jobs.CleanupSession(ctx, m.handlers.RemoveExpiredSessions)
}

func (m *AuthModule) Middleware() echo.MiddlewareFunc {
	return api.AuthMiddleware(m.handlers.VerifySession)
}
