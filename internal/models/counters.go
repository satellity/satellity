package models

import (
	"context"
	"database/sql"
	"godiscourse/internal/session"
)

func usersCount(ctx context.Context, tx *sql.Tx) (int64, error) {
	var count int64
	err := tx.QueryRowContext(ctx, "SELECT count(*) FROM users").Scan(&count)
	return count, err
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

func topicsCount(ctx context.Context, tx *sql.Tx) (int64, error) {
	var count int64
	err := tx.QueryRowContext(ctx, "SELECT count(*) FROM topics").Scan(&count)
	return count, err
}

func commentsCount(ctx context.Context, tx *sql.Tx) (int64, error) {
	var count int64
	err := tx.QueryRowContext(ctx, "SELECT count(*) FROM comments").Scan(&count)
	return count, err
}
