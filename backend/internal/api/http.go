package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v5"
)

type HttpRouter struct {
	router http.Handler
	api    huma.API
}

func NewHttpRouter() *HttpRouter {
	router := echo.New()

	humaConfig := huma.DefaultConfig("File Ctrl API", "v1")
	humaConfig.DocsRenderer = huma.DocsRendererScalar
	humaConfig.CreateHooks = nil

	api := humaecho.New(router, humaConfig)
	apiGroup := huma.NewGroup(api, "/api")

	registerFrontend(router)

	return &HttpRouter{
		router: router,
		api:    apiGroup,
	}
}

func (r *HttpRouter) RegisterRoutes(register func(api huma.API)) {
	register(r.api)
}

func (r *HttpRouter) Run(port int) error {
	server := http.Server{
		Addr:         fmt.Sprintf(":%v", port),
		Handler:      r.router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	slog.Info("Server started", slog.Int("PORT", port))

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start server: %w", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		_ = server.Close()
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}

	log.Println("HTTP server stopped")

	return nil
}
