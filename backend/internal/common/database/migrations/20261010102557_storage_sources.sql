-- +goose Up
CREATE TABLE IF NOT EXISTS storage_sources (
  path TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS storage_sources;
