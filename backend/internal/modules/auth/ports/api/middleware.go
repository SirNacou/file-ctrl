package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"uuid"

	"github.com/labstack/echo/v5"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/app"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/domain"
)

type (
	userContextKey struct{}
	UserContext    struct {
		ID       uuid.UUID
		Username string
		Role     domain.Role
	}
)

func GetUserFromContext(ctx context.Context) *UserContext {
	user, ok := ctx.Value(userContextKey{}).(*UserContext)
	if !ok {
		return nil
	}
	return user
}

var publicPaths = map[string]struct{}{
	"/api/auth/status": {},
	"/api/auth/login":  {},
	"/api/auth/setup":  {},
}

func AuthMiddleware(handler *app.VerifySessionHandler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			r := c.Request()
			path := r.URL.Path

			// Only protect /api/*; frontend assets pass through.
			if !strings.HasPrefix(path, "/api/") {
				return next(c)
			}
			if _, ok := publicPaths[strings.TrimSuffix(path, "/")]; ok {
				return next(c)
			}

			cookie, err := c.Cookie(SessionCookieName)
			if err != nil || cookie.Value == "" {
				return writeProblem(c, http.StatusUnauthorized, "Authentication required")
			}

			res, err := handler.Handle(
				r.Context(),
				&app.VerifySessionRequest{SessionToken: cookie.Value},
			)
			switch {
			case err == nil:
			case errors.Is(err, domain.ErrSessionNotFound):
				return writeProblem(c, http.StatusUnauthorized, "Session expired or invalid")
			case errors.Is(err, domain.ErrUserNotFound):
				return writeProblem(c, http.StatusUnauthorized, "User account no longer exists")
			default:
				slog.ErrorContext(r.Context(), "verify session failed", slog.Any("err", err))
				return writeProblem(c, http.StatusInternalServerError, "Internal error")
			}

			ctx := context.WithValue(r.Context(), userContextKey{}, &UserContext{
				ID:       res.UserID,
				Username: res.Username,
				Role:     res.Role,
			})
			c.SetRequest(r.WithContext(ctx))

			return next(c)
		}
	}
}

func writeProblem(c *echo.Context, status int, detail string) error {
	body, err := json.Marshal(map[string]any{
		"title":  http.StatusText(status),
		"status": status,
		"detail": detail,
	})
	if err != nil {
		return err
	}
	return c.Blob(status, "application/problem+json", body)
}
