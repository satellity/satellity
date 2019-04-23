package models

import (
	"context"
	"database/sql"
)

const (
	minCommentBodySize = 6
)

const commentsDDL = `
CREATE TABLE IF NOT EXISTS comments (
	comment_id            VARCHAR(36) PRIMARY KEY,
	body                  TEXT NOT NULL,
  topic_id              VARCHAR(36) NOT NULL REFERENCES topics ON DELETE CASCADE,
	user_id               VARCHAR(36) NOT NULL REFERENCES users ON DELETE CASCADE,
	score                 INTEGER NOT NULL DEFAULT 0,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE INDEX ON comments (topic_id, created_at);
CREATE INDEX ON comments (user_id, created_at);
CREATE INDEX ON comments (score DESC, created_at);
`

func commentsCount(ctx context.Context, tx *sql.Tx) (int64, error) {
	var count int64
	err := tx.QueryRowContext(ctx, "SELECT count(*) FROM comments").Scan(&count)
	return count, err
}
