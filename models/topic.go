package models

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/godiscourse/godiscourse/session"
	"github.com/godiscourse/godiscourse/uuid"
)

// Topic related CONST
const (
	MinimumTitleSize = 3
	MinimumBodySize  = 3
)

const topicsDDL = `
CREATE TABLE IF NOT EXISTS topics (
	topic_id              VARCHAR(36) PRIMARY KEY,
	title                 VARCHAR(512) NOT NULL,
	body                  TEXT NOT NULL,
	comments_count        INTEGER NOT NULL DEFAULT 0,
	category_id           VARCHAR(36) NOT NULL,
	user_id               VARCHAR(36) NOT NULL REFERENCES users ON DELETE CASCADE,
	score                 INTEGER NOT NULL DEFAULT 0,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX ON topics (user_id);
CREATE INDEX ON topics (category_id);
CREATE INDEX ON topics (created_at DESC);
CREATE INDEX ON topics (user_id, created_at DESC);
CREATE INDEX ON topics (score DESC, created_at DESC);
`

var topicCols = []string{"topic_id", "title", "body", "comments_count", "category_id", "user_id", "score", "created_at", "updated_at"}

// Topic is what use talking about
type Topic struct {
	TopicID       string    `sql:"topic_id,pk"`
	Title         string    `sql:"title"`
	Body          string    `sql:"body"`
	CommentsCount int       `sql:"comments_count,notnull"`
	CategoryID    string    `sql:"category_id"`
	UserID        string    `sql:"user_id"`
	Score         int       `sql:"score,notnull"`
	CreatedAt     time.Time `sql:"created_at"`
	UpdatedAt     time.Time `sql:"updated_at"`
}

//CreateTopic create a new Topic
func (user *User) CreateTopic(ctx context.Context, title, body, categoryID string) (*Topic, error) {
	title = strings.TrimSpace(title)
	body = strings.TrimSpace(body)
	if len(title) < MinimumTitleSize {
		return nil, session.BadDataError(ctx)
	}

	topic := &Topic{
		TopicID: uuid.NewV4().String(),
		Title:   title,
		Body:    body,
		UserID:  user.UserID,
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

// UpdateTopic update a Topic by ID, TODO maybe update categoryID also
func (user *User) UpdateTopic(ctx context.Context, id, title, body string) (*Topic, error) {
	title, body = strings.TrimSpace(title), strings.TrimSpace(body)
	if len(title) < MinimumTitleSize && len(body) < MinimumBodySize {
		return nil, session.BadDataError(ctx)
	}

	var topic *Topic
	err := session.Database(ctx).RunInTransaction(func(tx *pg.Tx) error {
		var err error
		topic, err = findTopic(ctx, tx, id)
		if err != nil {
			return err
		} else if topic == nil {
			return session.NotFoundError(ctx)
		} else if topic.UserID != user.UserID && !user.isAdmin() {
			return session.AuthorizationError(ctx)
		}
		if title != "" {
			topic.Title = title
		}
		if body != "" {
			topic.Body = body
		}
		return tx.Update(topic)
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
		return nil, err
	}
	return topic, nil
}

// ReadTopics read all topics, parameters: offset default time.Now()
func ReadTopics(ctx context.Context, offset time.Time) ([]*Topic, error) {
	if offset.IsZero() {
		offset = time.Now()
	}
	var topics []*Topic
	if _, err := session.Database(ctx).Query(&topics, "SELECT * FROM topics WHERE created_at<? ORDER BY created_at DESC LIMIT 50", offset); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topics, nil
}

// ReadTopics read user's topics, parameters: offset default time.Now()
func (user *User) ReadTopics(ctx context.Context, offset time.Time) ([]*Topic, error) {
	if offset.IsZero() {
		offset = time.Now()
	}
	var topics []*Topic
	if _, err := session.Database(ctx).Query(&topics, "SELECT * FROM topics WHERE user_id=? AND created_at<? ORDER BY created_at DESC LIMIT 50", user.UserID, offset); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topics, nil
}

// BeforeInsert hook insert
func (t *Topic) BeforeInsert(db orm.DB) error {
	t.CreatedAt = time.Now()
	t.UpdatedAt = t.CreatedAt
	return nil
}

// BeforeUpdate hook update
func (t *Topic) BeforeUpdate(db orm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}
