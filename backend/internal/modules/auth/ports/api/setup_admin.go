package api

import (
	"context"
	"errors"
	"net/http"
	"uuid"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/app"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/domain"
)

type SetupAdminRequest struct {
	Body SetupAdminRequestBody
}
type SetupAdminRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type SetupAdminResponse struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
	Body      SetupAdminResponseBody
}

type SetupAdminResponseBody struct {
	ID       uuid.UUID   `json:"user_id"`
	Username string      `json:"username"`
	Role     domain.Role `json:"role"`
}

func SetupAdmin(api huma.API, h *app.SetupAdminHandler, cookies *SessionCookies) {
	huma.Register(api, huma.Operation{
		OperationID: "setup-initial-admin",
		Method:      http.MethodPost,
		Path:        "/setup",
		Summary:     "Create first administrator account on first run",
	}, func(ctx context.Context, input *SetupAdminRequest) (*SetupAdminResponse, error) {
		res, err := h.Handle(ctx, &app.SetupAdminRequest{
			Username: input.Body.Username,
			Password: input.Body.Password,
		})
		if err != nil {
			if errors.Is(err, domain.ErrSetupAlreadyDone) {
				return nil, huma.Error409Conflict("Setup has already been completed", err)
			}
			return nil, huma.Error400BadRequest(err.Error(), err)
		}

		cookie := cookies.New(res.RawToken, res.ExpiresAt)

		return &SetupAdminResponse{
			SetCookie: *cookie,
			Body: SetupAdminResponseBody{
				ID:       res.User.ID,
				Username: res.User.Username,
				Role:     res.User.Role,
			},
		}, nil
	})
}
