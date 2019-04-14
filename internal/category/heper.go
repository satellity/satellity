package category

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"godiscourse/internal/durable"
	"godiscourse/internal/session"
)

var categoryColumns = []string{"category_id", "name", "alias", "description", "topics_count", "last_topic_id", "position", "created_at", "updated_at"}

func (m *Model) values() []interface{} {
	return []interface{}{m.CategoryID, m.Name, m.Alias, m.Description, m.TopicsCount, m.LastTopicID, m.Position, m.CreatedAt, m.UpdatedAt}
}

func categoryCount(ctx context.Context, tx *sql.Tx) (int64, error) {
	var count int64
	row := tx.QueryRowContext(ctx, "SELECT count(*) FROM categories")
	err := row.Scan(&count)
	if err != nil {
		return 0, session.TransactionError(ctx, err)
	}
	return count, nil
}

func categoryFromRows(row durable.Row) (m *Model, err error) {
	err = row.Scan(&m.CategoryID, &m.Name, &m.Alias, &m.Description, &m.TopicsCount, &m.LastTopicID, &m.Position, &m.CreatedAt, &m.UpdatedAt)
	return
}

func readCategories(ctx context.Context, tx *sql.Tx) ([]*Model, error) {
	rows, err := tx.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM categories ORDER BY position LIMIT 500", strings.Join(categoryColumns, ",")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Model
	for rows.Next() {
		category, err := categoryFromRows(rows)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}
