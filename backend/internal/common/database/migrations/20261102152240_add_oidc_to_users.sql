-- +goose Up
ALTER TABLE users
ADD COLUMN provider TEXT NOT NULL DEFAULT 'local';

ALTER TABLE users
ADD COLUMN provider_id TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_provider_identity ON users (provider, provider_id)
WHERE
  provider_id NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_users_provider_identity;
ALTER TABLE users DROP COLUMN provider_id;
ALTER TABLE users DROP COLUMN provider;
