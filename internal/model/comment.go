package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"godiscourse/internal/durable"

	"github.com/gofrs/uuid"
)

type Comment struct {
	CommentID string    `sql:"comment_id,pk"`
	Body      string    `sql:"body"`
	TopicID   string    `sql:"topic_id"`
	UserID    string    `sql:"user_id"`
	Score     int       `sql:"score,notnull"`
	CreatedAt time.Time `sql:"created_at"`
	UpdatedAt time.Time `sql:"updated_at"`
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
