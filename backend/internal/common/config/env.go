package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/caarlos0/env/v11"
)

var (
	ErrInvalidPort     = errors.New("port must not be empty")
	ErrDBInsideStorage = errors.New("database file cannot be located inside the storage root")
)

type Env struct {
	AppEnv               string `env:"APP_ENV"       envDefault:"development"`
	DatabasePath         string `env:"DATABASE_PATH" envDefault:"../data/metadata.db"`
	StorageDir           string `env:"STORAGE_DIR"   envDefault:"../storage"`
	CookieSecureOverride *bool  `env:"COOKIE_SECURE"`
}

func LoadEnv() (*Env, error) {
	cfg := new(Env)
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Env) IsProd() bool {
	return c.AppEnv == "production"
}

func (c *Env) IsDev() bool {
	return c.AppEnv == "development"
}

func (c *Env) Port() int {
	return 8080
}

// Validate checks pure in-memory constraints and normalizes paths.
// No filesystem operations (os.MkdirAll, os.Stat) are executed here.
func (c *Env) Validate() error {
	if c.Port() < 0 || c.Port() > 65535 {
		return ErrInvalidPort
	}

	// Normalize storage root to an absolute path
	absStorage, err := filepath.Abs(filepath.Clean(c.StorageDir))
	if err != nil {
		return fmt.Errorf("invalid storage root: %w", err)
	}

	c.StorageDir = absStorage

	// Normalize database path to an absolute path
	absDB, err := filepath.Abs(filepath.Clean(c.DatabasePath))
	if err != nil {
		return fmt.Errorf("invalid database path: %w", err)
	}

	c.DatabasePath = absDB

	// Assert SQLite database is not located inside the storage directory
	rel, err := filepath.Rel(c.StorageDir, c.DatabasePath)
	if err == nil && !strings.HasPrefix(rel, "..") {
		return ErrDBInsideStorage
	}

	return nil
}

func (c *Env) CookieSecure() bool {
	if c.CookieSecureOverride != nil {
		return *c.CookieSecureOverride
	}
	return c.IsProd()
}
