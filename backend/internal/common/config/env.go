package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	AppEnv string `env:"APP_ENV" envDefault="development"`
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

func (c *Config) StorageRoot() string {
	if c.IsProd() {
		return "/storage"
	} else {
		// Local development inside workspace
		return "../storage"
	}
}

func (c *Config) DbPath() string {
	if c.IsProd() {
		return "/data/metadata.db"
	} else {
		return "../data/metadata.db"
	}
}
