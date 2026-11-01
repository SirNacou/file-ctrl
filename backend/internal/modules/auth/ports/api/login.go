package api

import (
	"context"
	"net/http"
	"uuid"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/app"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/domain"
)

type LoginRequest struct {
	Body LoginRequestBody
}
type LoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
	Body      LoginResponseBody
}
type LoginResponseBody struct {
	ID       uuid.UUID   `json:"user_id"`
	Username string      `json:"username"`
	Role     domain.Role `json:"role"`
}

func Login(api huma.API, h *app.LoginHandler, cookies *SessionCookies) {
	huma.Register(api, huma.Operation{
		OperationID: "auth-login",
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "Authenticate with username and password",
	}, func(ctx context.Context, input *LoginRequest) (*LoginResponse, error) {
		res, err := h.Handle(ctx, &app.LoginRequest{
			Username: input.Body.Username,
			Password: input.Body.Password,
		})
		if err != nil {
			return nil, huma.Error401Unauthorized("Invalid username or password", err)
		}

		cookie := cookies.New(res.RawToken, res.ExpiresAt)

		return &LoginResponse{
			Body: LoginResponseBody{
				ID:       res.User.ID,
				Username: res.User.Username,
				Role:     res.User.Role,
			},
			SetCookie: *cookie,
		}, nil
	})
}
