package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	AppEnv     string `env:"APP_ENV" envDefault:"development"`
	DataDir    string `env:"DATA_DIR" envDefault:"../data"`
	StorageDir string `env:"STORAGE_DIR" envDefault:"../storage"`
}

func LoadEnv() (*Config, error) {
	cfg := new(Config)
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	return cfg, nil
}

func (c *Config) IsProd() bool {
	return c.AppEnv == "production"
}

func (c *Config) IsDev() bool {
	return c.AppEnv == "development"
}

func (c *Config) Port() uint16 {
	return 8080
}
