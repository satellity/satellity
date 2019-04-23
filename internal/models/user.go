package models

import (
	"context"
	"database/sql"
)

const usersDDL = `
CREATE TABLE IF NOT EXISTS users (
	user_id               VARCHAR(36) PRIMARY KEY,
	email                 VARCHAR(512),
	username              VARCHAR(64) NOT NULL CHECK (username ~* '^[a-z0-9][a-z0-9_]{3,63}$'),
	nickname              VARCHAR(64) NOT NULL DEFAULT '',
	biography             VARCHAR(2048) NOT NULL DEFAULT '',
	encrypted_password    VARCHAR(1024),
	github_id             VARCHAR(1024) UNIQUE,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX ON users ((LOWER(email)));
CREATE UNIQUE INDEX ON users ((LOWER(username)));
CREATE INDEX ON users (created_at);
`

// todo: move to statistics
func usersCount(ctx context.Context, tx *sql.Tx) (int64, error) {
	var count int64
	err := tx.QueryRowContext(ctx, "SELECT count(*) FROM users").Scan(&count)
	return count, err
}
