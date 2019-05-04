package topic

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"godiscourse/internal/durable"

	"github.com/gofrs/uuid"
)

var topicColumns = []string{"topic_id", "short_id", "title", "body", "comments_count", "category_id", "user_id", "score", "created_at", "updated_at"}

func topicFromRows(row durable.Row) (*Model, error) {
	var m Model
	err := row.Scan(&m.TopicID, &m.ShortID, &m.Title, &m.Body, &m.CommentsCount, &m.CategoryID, &m.UserID, &m.Score, &m.CreatedAt, &m.UpdatedAt)
	return &m, err
}

func findTopic(ctx context.Context, tx *sql.Tx, id string) (*Model, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}
	row := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM topics WHERE topic_id=$1", strings.Join(topicColumns, ",")), id)
	t, err := topicFromRows(row)
	if sql.ErrNoRows == err {
		return nil, nil
	}
	return t, err
}

func findTopicByShortID(ctx context.Context, tx *sql.Tx, id string) (*Model, error) {
	row := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM topics WHERE short_id=$1", strings.Join(topicColumns, ",")), id)
	t, err := topicFromRows(row)
	if sql.ErrNoRows == err {
		return nil, nil
	}
	return t, err
}

func topicsCountByCategory(ctx context.Context, tx *sql.Tx, id string) (int64, error) {
	var count int64
	err := tx.QueryRowContext(ctx, "SELECT count(*) FROM topics WHERE category_id=$1", id).Scan(&count)
	return count, err
}
