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

const topics_DDL = `
CREATE TABLE IF NOT EXISTS topics (
	topic_id              VARCHAR(36) PRIMARY KEY,
	title                 VARCHAR(512) NOT NULL,
	body                  TEXT NOT NULL,
	category_id           VARCHAR(36) NOT NULL,
	user_id               VARCHAR(36) NOT NULL,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX ON topics (category_id);
CREATE INDEX ON topics (user_id);
`

var topicCols = []string{"topic_id", "title", "body", "category_id", "user_id", "created_at", "updated_at"}

type Topic struct {
	TopicId    string    `sql:"topic_id"`
	Title      string    `sql:"title"`
	Body       string    `sql:"body"`
	CategoryId string    `sql:"category_id"`
	UserId     string    `sql:"user_id"`
	CreatedAt  time.Time `sql:"created_at"`
	UpdatedAt  time.Time `sql:"updated_at"`
}

func (user *User) CreateTopic(ctx context.Context, title, body, categoryId string) (*Topic, error) {
	title := strings.TrimSpace(title)
	body := strings.TrimSpace(body)
	if len(title) < 1 {
		return nil, session.BadDataError(ctx)
	}

	t := time.Now()
	topic := &Topic{
		TopicId:   uuid.NewV4().String(),
		Title:     title,
		Body:      body,
		UserId:    user.UserId,
		CreatedAt: t,
		UpdatedAt: t,
	}
	err := session.Database(ctx).RunInTransaction(func(tx *pg.Tx) error {
		category, err := findCategory(ctx, tx, categoryId)
		if err != nil {
			return err
		}
		if category == nil {
			return session.BadDataError(ctx)
		}
		topic.CategoryId = category.CategoryId
		category.LastTopicId = sql.NullString{topic.TopicId, true}
		category.TopicsCount += 1
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

func ReadTopic(ctx context.Context, id string) (*Topic, error) {
	topic := &Topic{TopicId: id}
	if err := session.Database(ctx).Model(topic).Column(topicCols...).WherePK().Select(); err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topic, nil
}
