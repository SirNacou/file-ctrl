package database

import (
	"context"
	"fmt"

	"github.com/stephenafamo/bob"
	_ "modernc.org/sqlite"
)

func NewDb(ctx context.Context, dbPath string) (*bob.DB, error) {
	db, err := bob.Open("sqlite", fmt.Sprintf("file:%s?_journal=WAL&_timeout=5000", dbPath))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return &db, nil
}
