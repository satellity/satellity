package engine

import (
	"context"
	"database/sql"

	"godiscourse/internal/session"
)

func categoryCount(ctx context.Context, tx *sql.Tx) (int64, error) {
	var count int64
	row := tx.QueryRowContext(ctx, "SELECT count(*) FROM categories")
	err := row.Scan(&count)
	if err != nil {
		return 0, session.TransactionError(ctx, err)
	}
	return count, nil
}

func topicsCountByCategory(ctx context.Context, tx *sql.Tx, id string) (int64, error) {
	var count int64
	err := tx.QueryRowContext(ctx, "SELECT count(*) FROM topics WHERE category_id=$1", id).Scan(&count)
	return count, err
}

func commentsCountByTopic(ctx context.Context, tx *sql.Tx, id string) (int64, error) {
	var count int64
	err := tx.QueryRowContext(ctx, "SELECT count(*) FROM comments WHERE topic_id=$1", id).Scan(&count)
	return count, err
}
