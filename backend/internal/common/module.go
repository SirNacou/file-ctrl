package common

import (
	"context"

	"github.com/sirnacou/file-ctrl/backend/internal/common/config"
	"github.com/sirnacou/file-ctrl/backend/internal/common/database"
	"github.com/stephenafamo/bob"
)

type CommonModule struct {
	DB *bob.DB
}

func NewCommonModule(ctx context.Context, cfg *config.Config) (*CommonModule, error) {
	db, err := database.NewDb(ctx, cfg.DbPath())
	if err != nil {
		return nil, err
	}

	err = database.RunMigrations(db.DB)
	if err != nil {
		return nil, err
	}
	return &CommonModule{
		DB: db,
	}, nil
}
