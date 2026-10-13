-- +goose Up
CREATE TABLE IF NOT EXISTS storage_sources (
  id TEXT PRIMARY KEY, -- UUID or slug (e.g. "photos")
  name TEXT NOT NULL, -- UI Display Name (e.g. "Family Photos")
  subpath TEXT NOT NULL UNIQUE, -- Directory name under /storage (e.g. "photos")
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS storage_sources;
