package main

import (
	"context"
	"database/sql"
	"fmt"
	"godiscourse/internal/durable"
	"godiscourse/internal/session"
	"time"

	hashids "github.com/speps/go-hashids"
)

func generateShortID(table string, t time.Time) (string, error) {
	hd := hashids.NewData()
	hd.MinLength = 5
	h, _ := hashids.NewWithData(hd)
	return h.EncodeInt64([]int64{t.UnixNano()})
}

// MigrateTopics should be deleted after task TODO
func MigrateTopics(db *durable.Database, offset time.Time, limit int64) (int64, time.Time, error) {
	ctx := context.Background()
	if offset.IsZero() {
		offset = time.Now()
	}

	last := offset
	var count int64
	set := make(map[string]string)
	err := db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		query := "SELECT topic_id,short_id,created_at FROM topics WHERE created_at<$1 ORDER BY created_at DESC LIMIT $2"
		rows, err := tx.QueryContext(ctx, query, offset, limit)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var topicID string
			var shortID sql.NullString
			var t time.Time
			err = rows.Scan(&topicID, &shortID, &t)
			if err != nil {
				return err
			}
			count++
			last = t
			if shortID.Valid {
				continue
			}
			id, _ := generateShortID("topics", last)
			set[topicID] = id
		}
		if err := rows.Err(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, offset, session.TransactionError(ctx, err)
	}
	for k, v := range set {
		query := fmt.Sprintf("UPDATE topics SET short_id='%s' WHERE topic_id='%s'", v, k)
		_, err = db.ExecContext(ctx, query)
		if err != nil {
			return 0, offset, session.TransactionError(ctx, err)
		}
	}
	return count, last, nil
}
