// internal/api/frontend.go
package api

import (
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/sirnacou/file-ctrl/backend/web"
)

func registerFrontend(e *echo.Echo) {
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       ".output/public",
		Index:      "_shell.html",
		Filesystem: web.FS,
		HTML5:      true,
		// Don't serve index.html for unknown API paths; let them 404 properly
		Skipper: func(c *echo.Context) bool {
			return strings.HasPrefix(c.Request().URL.Path, "/api")
		},
	}))
}
