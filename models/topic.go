package models

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/go-pg/pg"
	"github.com/godiscourse/godiscourse/session"
	"github.com/godiscourse/godiscourse/uuid"
)

const topicsDDL = `
CREATE TABLE IF NOT EXISTS topics (
	topic_id              VARCHAR(36) PRIMARY KEY,
	title                 VARCHAR(512) NOT NULL,
	body                  TEXT NOT NULL,
	category_id           VARCHAR(36) NOT NULL,
	user_id               VARCHAR(36) NOT NULL,
	score                 INTEGER NOT NULL DEFAULT 0,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX ON topics (user_id);
CREATE INDEX ON topics (category_id);
CREATE INDEX ON topics (created_at DESC);
CREATE INDEX ON topics (score, created_at DESC);
`

var topicCols = []string{"topic_id", "title", "body", "category_id", "user_id", "score", "created_at", "updated_at"}

// Topic is what use talking about
type Topic struct {
	TopicID    string    `sql:"topic_id"`
	Title      string    `sql:"title"`
	Body       string    `sql:"body"`
	CategoryID string    `sql:"category_id"`
	UserID     string    `sql:"user_id"`
	Score      int       `sql:"score"`
	CreatedAt  time.Time `sql:"created_at"`
	UpdatedAt  time.Time `sql:"updated_at"`
}

//CreateTopic create a new Topic
func (user *User) CreateTopic(ctx context.Context, title, body, categoryID string) (*Topic, error) {
	title = strings.TrimSpace(title)
	body = strings.TrimSpace(body)
	if len(title) < 1 {
		return nil, session.BadDataError(ctx)
	}

	t := time.Now()
	topic := &Topic{
		TopicID:   uuid.NewV4().String(),
		Title:     title,
		Body:      body,
		UserID:    user.UserID,
		CreatedAt: t,
		UpdatedAt: t,
	}
	err := session.Database(ctx).RunInTransaction(func(tx *pg.Tx) error {
		category, err := findCategory(ctx, tx, categoryID)
		if err != nil {
			return err
		}
		if category == nil {
			return session.BadDataError(ctx)
		}
		topic.CategoryID = category.CategoryID
		category.LastTopicID = sql.NullString{String: topic.TopicID, Valid: true}
		category.TopicsCount++
		if err := tx.Insert(topic); err != nil {
			return err
		}
		return tx.Update(category)
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return nil, session.TransactionError(ctx, err)
	}
	return topic, nil
}

//ReadTopic read a topic by ID
func ReadTopic(ctx context.Context, id string) (*Topic, error) {
	var topic *Topic
	err := session.Database(ctx).RunInTransaction(func(tx *pg.Tx) error {
		var err error
		topic, err = findTopic(ctx, tx, id)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topic, nil
}

func findTopic(ctx context.Context, tx *pg.Tx, id string) (*Topic, error) {
	topic := &Topic{TopicID: id}
	if err := tx.Model(topic).Column(topicCols...).WherePK().Select(); err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topic, nil
}
