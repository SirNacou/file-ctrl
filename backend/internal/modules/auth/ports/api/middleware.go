package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"uuid"

	"github.com/danielgtaylor/huma/v2"
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

func AuthMiddleware(
	api huma.API,
	handler *app.VerifySessionHandler,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			// 1. Only protect /api/* routes (let frontend static assets and /docs pass)
			if !strings.HasPrefix(path, "/api/") {
				next.ServeHTTP(w, r)
				return
			}

			// 2. Allow public auth routes (status, login, setup)
			if strings.HasPrefix(path, "/api/auth/status") ||
				strings.HasPrefix(path, "/api/auth/login") ||
				strings.HasPrefix(path, "/api/auth/setup") {
				next.ServeHTTP(w, r)
				return
			}

			// 2. Read session cookie
			cookie, err := r.Cookie(SessionCookieName)
			if err != nil || cookie.Value == "" {
				writeUnauthorized(w, "Authentication required")
				return
			}

			res, err := handler.Handle(
				r.Context(),
				&app.VerifySessionRequest{SessionToken: cookie.Value},
			)
			if err != nil {
				if errors.Is(err, domain.ErrSessionNotFound) {
					writeUnauthorized(w, "Session expired or invalid")
				} else if errors.Is(err, domain.ErrUserNotFound) {
					writeUnauthorized(w, "User account no longer exists")
				}
				return
			}

			// 5. Inject into context and proceed
			ctx := context.WithValue(r.Context(), userContextKey{}, UserContext{
				ID:       res.UserID,
				Username: res.Username,
				Role:     res.Role,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeUnauthorized(w http.ResponseWriter, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"title":  "Unauthorized",
		"status": http.StatusUnauthorized,
		"detail": detail,
	})
}
