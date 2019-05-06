package models

import (
	"context"
	"database/sql"
	"fmt"
	"godiscourse/internal/durable"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	hashids "github.com/speps/go-hashids"
)

type Topic struct {
	TopicID       string
	ShortID       string
	Title         string
	Body          string
	CommentsCount int64
	CategoryID    string
	UserID        string
	Score         int
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Category Category
	User     User
}

type TopicInfo struct {
	Title      string
	Body       string
	CategoryID string
}

var TopicColumns = []string{"topic_id", "short_id", "title", "body", "comments_count", "category_id", "user_id", "score", "created_at", "updated_at"}

func (t *Topic) Values() []interface{} {
	return []interface{}{t.TopicID, t.ShortID, t.Title, t.Body, t.CommentsCount, t.CategoryID, t.UserID, t.Score, t.CreatedAt, t.UpdatedAt}
}

func TopicFromRows(row durable.Row) (*Topic, error) {
	var t Topic
	err := row.Scan(&t.TopicID, &t.ShortID, &t.Title, &t.Body, &t.CommentsCount, &t.CategoryID, &t.UserID, &t.Score, &t.CreatedAt, &t.UpdatedAt)
	return &t, err
}

func GenerateShortID(table string, t time.Time) (string, error) {
	hd := hashids.NewData()
	hd.MinLength = 5
	h, _ := hashids.NewWithData(hd)
	return h.EncodeInt64([]int64{t.UnixNano()})
}

func FindTopic(ctx context.Context, tx *sql.Tx, id string) (*Topic, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}
	row := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM topics WHERE topic_id=$1", strings.Join(TopicColumns, ",")), id)
	t, err := TopicFromRows(row)
	if sql.ErrNoRows == err {
		return nil, nil
	}
	return t, err
}

func FindTopicByShortID(ctx context.Context, tx *sql.Tx, id string) (*Topic, error) {
	row := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM topics WHERE short_id=$1", strings.Join(TopicColumns, ",")), id)
	t, err := TopicFromRows(row)
	if sql.ErrNoRows == err {
		return nil, nil
	}
	return t, err
}

func LastTopic(ctx context.Context, categoryID string, tx *sql.Tx) (*Topic, error) {
	row := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM topics WHERE category_id=$1 LIMIT 1", strings.Join(TopicColumns, ",")), categoryID)
	t, err := TopicFromRows(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return t, err
}
