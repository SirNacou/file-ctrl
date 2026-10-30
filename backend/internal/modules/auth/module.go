package auth

import (
	"context"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/adapters"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/app"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/ports/jobs"
	"github.com/stephenafamo/bob"
)

type AuthModule struct {
	removeExpiredSessionsHandler *app.RemoveExpiredSessionsHandler
}

func NewAuthModule(db bob.Executor) *AuthModule {
	// userRepo := adapters.NewUserRepository(db)
	sessionRepo := adapters.NewSessionRepository(db)

	return &AuthModule{
		removeExpiredSessionsHandler: app.NewRemoveExpiredSessionsHandler(sessionRepo),
	}
}

func (m *AuthModule) RegisterJobs(ctx context.Context) {
	jobs.CleanupSession(ctx, m.removeExpiredSessionsHandler)
}
