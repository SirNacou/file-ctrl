package config

import (
	"path/filepath"
	"testing"
)

func TestLoadEnvUsesPortEnvironmentVariable(t *testing.T) {
	t.Setenv("PORT", "18081")
	t.Setenv("APP_ENV", "development")
	t.Setenv("OIDC_ENABLED", "false")

	testDir := t.TempDir()
	t.Setenv("DATABASE_PATH", filepath.Join(testDir, "metadata.db"))
	t.Setenv("STORAGE_DIR", filepath.Join(testDir, "storage"))
	t.Setenv("UPLOAD_DIR", filepath.Join(testDir, "uploads"))

	cfg, err := LoadEnv()
	if err != nil {
		t.Fatalf("LoadEnv() error = %v", err)
	}
	if got, want := cfg.Port(), 18081; got != want {
		t.Errorf("Port() = %d, want %d", got, want)
	}
}
