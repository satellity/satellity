package models

import (
	"context"
	"database/sql"
	"fmt"
	"godiscourse/internal/durable"
	"godiscourse/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	hashids "github.com/speps/go-hashids"
)

// Topic related CONST
const (
	minTitleSize = 3
	LIMIT        = 50
)

const topicsDDL = `
CREATE TABLE IF NOT EXISTS topics (
	topic_id              VARCHAR(36) PRIMARY KEY,
	short_id              VARCHAR(255) NOT NULL,
	title                 VARCHAR(512) NOT NULL,
	body                  TEXT NOT NULL,
	comments_count        INTEGER NOT NULL DEFAULT 0,
	category_id           VARCHAR(36) NOT NULL,
	user_id               VARCHAR(36) NOT NULL REFERENCES users ON DELETE CASCADE,
	score                 INTEGER NOT NULL DEFAULT 0,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX ON topics (short_id);
CREATE INDEX ON topics (created_at DESC);
CREATE INDEX ON topics (user_id, created_at DESC);
CREATE INDEX ON topics (category_id, created_at DESC);
CREATE INDEX ON topics (score DESC, created_at DESC);
`

var topicColumns = []string{"topic_id", "short_id", "title", "body", "comments_count", "category_id", "user_id", "score", "created_at", "updated_at"}

func (t *Topic) values() []interface{} {
	return []interface{}{t.TopicID, t.ShortID, t.Title, t.Body, t.CommentsCount, t.CategoryID, t.UserID, t.Score, t.CreatedAt, t.UpdatedAt}
}

// Topic is what use talking about
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

	User     *User
	Category *Category
}

func findTopic(ctx context.Context, tx *sql.Tx, id string) (*Topic, error) {
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

func topicsCount(ctx context.Context, tx *sql.Tx) (int64, error) {
	var count int64
	err := tx.QueryRowContext(ctx, "SELECT count(*) FROM topics").Scan(&count)
	return count, err
}

func topicFromRows(row durable.Row) (*Topic, error) {
	var t Topic
	err := row.Scan(&t.TopicID, &t.ShortID, &t.Title, &t.Body, &t.CommentsCount, &t.CategoryID, &t.UserID, &t.Score, &t.CreatedAt, &t.UpdatedAt)
	return &t, err
}

func generateShortID(table string, t time.Time) (string, error) {
	hd := hashids.NewData()
	hd.MinLength = 5
	h, _ := hashids.NewWithData(hd)
	return h.EncodeInt64([]int64{t.UnixNano()})
}

// MigrateTopics should be deleted after task TODO
func MigrateTopics(mctx *Context, offset time.Time, limit int64) (int64, time.Time, error) {
	ctx := mctx.context
	if offset.IsZero() {
		offset = time.Now()
	}

	last := offset
	var count int64
	set := make(map[string]string)
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
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
		_, err = mctx.database.ExecContext(ctx, query)
		if err != nil {
			return 0, offset, session.TransactionError(ctx, err)
		}
	}
	return count, last, nil
}
