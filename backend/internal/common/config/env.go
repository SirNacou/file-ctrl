package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/caarlos0/env/v11"
)

var (
	ErrInvalidPort         = errors.New("port must not be empty")
	ErrDBInsideStorage     = errors.New("database file cannot be located inside the storage root")
	ErrUploadInsideStorage = errors.New("upload directory cannot be located inside the storage root")
)

type Env struct {
	AppEnv               string `env:"APP_ENV"       envDefault:"development"`
	DatabasePath         string `env:"DATABASE_PATH" envDefault:"../data/metadata.db"`
	StorageDir           string `env:"STORAGE_DIR"   envDefault:"../storage"`
	UploadDir            string `env:"UPLOAD_DIR"    envDefault:"../data/uploads"`
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

	// Normalize upload staging root to an absolute path
	absUpload, err := filepath.Abs(filepath.Clean(c.UploadDir))
	if err != nil {
		return fmt.Errorf("invalid upload root: %w", err)
	}
	c.UploadDir = absUpload

	// Normalize database path to an absolute path
	absDB, err := filepath.Abs(filepath.Clean(c.DatabasePath))
	if err != nil {
		return fmt.Errorf("invalid database path: %w", err)
	}

	c.DatabasePath = absDB

	// Assert SQLite database is not located inside the storage directory
	if pathIsWithin(c.StorageDir, c.DatabasePath) {
		return ErrDBInsideStorage
	}

	// Upload staging must stay outside the storage root.
	if pathIsWithin(c.StorageDir, c.UploadDir) {
		return ErrUploadInsideStorage
	}

	return nil
}

func pathIsWithin(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}

func (c *Env) CookieSecure() bool {
	if c.CookieSecureOverride != nil {
		return *c.CookieSecureOverride
	}
	return c.IsProd()
}
