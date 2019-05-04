package engine

import (
	"context"
	"database/sql"
)

func topicsCountByCategory(ctx context.Context, tx *sql.Tx, id string) (int64, error) {
	var count int64
	err := tx.QueryRowContext(ctx, "SELECT count(*) FROM topics WHERE category_id=$1", id).Scan(&count)
	return count, err
}