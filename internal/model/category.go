package model

import (
	"context"
	"database/sql"
	"fmt"
	"godiscourse/internal/durable"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

type Category struct {
	CategoryID  string
	Name        string
	Alias       string
	Description string
	TopicsCount int64
	LastTopicID sql.NullString
	Position    int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CategoryInfo struct {
	Name        string
	Alias       string
	Description string
	Position    int64
}

var CategoryColumns = []string{"category_id", "name", "alias", "description", "topics_count", "last_topic_id", "position", "created_at", "updated_at"}

func (c *Category) Values() []interface{} {
	return []interface{}{c.CategoryID, c.Name, c.Alias, c.Description, c.TopicsCount, c.LastTopicID, c.Position, c.CreatedAt, c.UpdatedAt}
}

func CategoryFromRows(row durable.Row) (*Category, error) {
	var c Category
	err := row.Scan(&c.CategoryID, &c.Name, &c.Alias, &c.Description, &c.TopicsCount, &c.LastTopicID, &c.Position, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func FindCategory(ctx context.Context, tx *sql.Tx, id string) (*Category, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}

	row := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM categories WHERE category_id=$1", strings.Join(CategoryColumns, ",")), id)
	cat, err := CategoryFromRows(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return cat, err
}
