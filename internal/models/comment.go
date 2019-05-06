package models

import (
	"context"
	"database/sql"
	"fmt"
	"godiscourse/internal/durable"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

const CommentsDDL = `
CREATE TABLE IF NOT EXISTS comments (
	comment_id            VARCHAR(36) PRIMARY KEY,
	body                  TEXT NOT NULL,
	topic_id              VARCHAR(36) NOT NULL REFERENCES topics ON DELETE CASCADE,
	user_id               VARCHAR(36) NOT NULL REFERENCES users ON DELETE CASCADE,
	score                 INTEGER NOT NULL DEFAULT 0,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS comments_topic_createdx ON comments (topic_id, created_at);
CREATE INDEX IF NOT EXISTS comments_user_createdx ON comments (user_id, created_at);
CREATE INDEX IF NOT EXISTS comments_score_createdx ON comments (score DESC, created_at);
`

const DropCommentsDDL = `DROP TABLE IF EXISTS comments;`

type Comment struct {
	CommentID string    `sql:"comment_id,pk"`
	Body      string    `sql:"body"`
	TopicID   string    `sql:"topic_id"`
	UserID    string    `sql:"user_id"`
	Score     int       `sql:"score,notnull"`
	CreatedAt time.Time `sql:"created_at"`
	UpdatedAt time.Time `sql:"updated_at"`
	User      User
}

type CommentInfo struct {
	CommentID string
	TopicID   string
	UserID    string
	Body      string
}

var CommentColumns = []string{"comment_id", "body", "topic_id", "user_id", "score", "created_at", "updated_at"}

func (c *Comment) Values() []interface{} {
	return []interface{}{c.CommentID, c.Body, c.TopicID, c.UserID, c.Score, c.CreatedAt, c.UpdatedAt}
}

func CommentFromRows(row durable.Row) (*Comment, error) {
	var c Comment
	err := row.Scan(&c.CommentID, &c.Body, &c.TopicID, &c.UserID, &c.Score, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func FindComment(ctx context.Context, tx *sql.Tx, id string) (*Comment, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}
	row := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM comments WHERE comment_id=$1", strings.Join(CommentColumns, ",")), id)
	c, err := CommentFromRows(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}
