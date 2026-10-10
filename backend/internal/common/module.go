package common

import (
	"context"
	"log/slog"
	"os"

	"github.com/sirnacou/file-ctrl/backend/internal/common/config"
	"github.com/sirnacou/file-ctrl/backend/internal/common/database"
	"github.com/stephenafamo/bob"
)

type CommonModule struct {
	Env    *config.Config
	Logger *slog.Logger
	DB     *bob.DB
}

func NewCommonModule(ctx context.Context) (*CommonModule, error) {
	env, err := config.LoadEnv()
	if err != nil {
		return nil, err
	}

	logger := newLogger(env.IsProd())
	db, err := newDB(ctx, env.DbPath())
	if err != nil {
		return nil, err
	}

	return &CommonModule{
		Env:    env,
		Logger: logger,
		DB:     db,
	}, nil
}

func newLogger(isProd bool) *slog.Logger {
	var logHandler slog.Handler
	if isProd {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	return slog.New(logHandler)
}

func newDB(ctx context.Context, dbPath string) (*bob.DB, error) {
	db, err := database.NewDb(ctx, dbPath)
	if err != nil {
		return nil, err
	}

	err = database.RunMigrations(db.DB)
	if err != nil {
		return nil, err
	}

	return db, nil
}
