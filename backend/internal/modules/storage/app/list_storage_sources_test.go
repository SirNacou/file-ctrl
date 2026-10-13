package app

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/sirnacou/file-ctrl/backend/internal/common/database"
	"github.com/stephenafamo/bob"
)

func TestListStorageSourcesHandler(t *testing.T) {
	db := newTestDB(t)

	createdAt := time.Date(2026, time.October, 13, 6, 33, 29, 294819074, time.UTC)

	_, err := db.ExecContext(t.Context(), `
		INSERT INTO storage_sources (id, name, subpath, created_at)
		VALUES (?, ?, ?, ?)
	`, "source-1", "Documents", "docs", createdAt)
	if err != nil {
		t.Fatalf("insert storage source: %v", err)
	}

	response, err := NewListStorageSourcesHandler(
		db,
	).Handle(t.Context(), &ListStorageSourcesRequest{})
	if err != nil {
		t.Fatalf("list storage sources: %v", err)
	}

	if len(response.Sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(response.Sources))
	}

	source := response.Sources[0]
	if source.ID != "source-1" || source.Name != "Documents" || source.Subpath != "docs" {
		t.Fatalf("unexpected source: %+v", source)
	}

	if !source.CreatedAt.Equal(createdAt) {
		t.Fatalf("expected created_at %s, got %s", createdAt, source.CreatedAt)
	}
}

func newTestDB(t *testing.T) *bob.DB {
	t.Helper()

	db, err := database.NewDb(
		t.Context(),
		filepath.Join(t.TempDir(), "metadata.db"),
	)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	t.Cleanup(func() {
		if err = db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	})

	_, err = db.ExecContext(t.Context(), `
		CREATE TABLE storage_sources (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			subpath TEXT NOT NULL UNIQUE,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("create table: %v", err)
	}

	return db
}
