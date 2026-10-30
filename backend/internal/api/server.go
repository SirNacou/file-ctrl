package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/sirnacou/file-ctrl/backend/internal/common"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	http    *http.Server
	common  *common.CommonModule
	storage *storage.StorageModule
	auth    *auth.AuthModule
}

func NewServer() (*Server, error) {
	cm, err := common.NewCommonModule()
	if err != nil {
		return nil, fmt.Errorf("create common module: %w", err)
	}
	slog.SetDefault(cm.Logger)

	e := newEcho()
	humaAPI := newHumaAPI(e)

	sm, err := storage.NewStorageModule(cm.Env.StorageDir, cm.DB)
	if err != nil {
		return nil, fmt.Errorf("create storage module: %w", err)
	}
	if err := sm.RegisterAPI(humaAPI, e); err != nil {
		return nil, fmt.Errorf("register storage api: %w", err)
	}

	registerFrontend(e)

	return &Server{
		http: &http.Server{
			Addr:              fmt.Sprintf(":%d", cm.Env.Port()),
			Handler:           e,
			ReadHeaderTimeout: 5 * time.Second,
			// no Read/WriteTimeout: uploads and downloads can be long
		},
		common:  cm,
		storage: sm,
		auth:    auth.NewAuthModule(cm.DB),
	}, nil
}

func newEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		ContentTypeNosniff: "nosniff",
	}))
	return e
}

func newHumaAPI(e *echo.Echo) huma.API {
	huma.DefaultArrayNullable = false // note: package-level global
	cfg := huma.DefaultConfig("File Ctrl API", "v1")
	cfg.DocsRenderer = huma.DocsRendererScalar
	cfg.CreateHooks = nil
	return huma.NewGroup(humaecho.New(e, cfg), "/api")
}

func (s *Server) Run(ctx context.Context) error {
	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		s.auth.RegisterJobs(gctx)
		return nil
	})

	g.Go(func() error {
		slog.Info("server started", slog.String("addr", s.http.Addr))
		if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		<-gctx.Done()
		slog.Info("shutting down")

		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
		defer cancel()

		if err := s.http.Shutdown(shutdownCtx); err != nil {
			_ = s.http.Close()
			return fmt.Errorf("shutdown http server: %w", err)
		}
		return nil
	})

	return g.Wait()
}
