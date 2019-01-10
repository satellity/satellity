package models

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/godiscourse/godiscourse/api/session"
	"github.com/satori/go.uuid"
)

// Topic related CONST
const (
	minTitleSize = 3
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

CREATE INDEX ON topics (created_at DESC);
CREATE INDEX ON topics (user_id, created_at DESC);
CREATE INDEX ON topics (category_id, created_at DESC);
CREATE INDEX ON topics (score DESC, created_at DESC);
`

var topicCols = []string{"topic_id", "title", "body", "comments_count", "category_id", "user_id", "score", "created_at", "updated_at"}

// Topic is what use talking about
type Topic struct {
	TopicID       string    `sql:"topic_id,pk"`
	Title         string    `sql:"title,notnull"`
	Body          string    `sql:"body,notnull"`
	CommentsCount int       `sql:"comments_count,notnull"`
	CategoryID    string    `sql:"category_id,notnull"`
	UserID        string    `sql:"user_id,notnull"`
	Score         int       `sql:"score,notnull"`
	CreatedAt     time.Time `sql:"created_at"`
	UpdatedAt     time.Time `sql:"updated_at"`

	User     *User     `sql:"-"`
	Category *Category `sql:"-"`
}

//CreateTopic create a new Topic
func (user *User) CreateTopic(ctx context.Context, title, body, categoryID string) (*Topic, error) {
	title = strings.TrimSpace(title)
	body = strings.TrimSpace(body)
	if len(title) < minTitleSize {
		return nil, session.BadDataError(ctx)
	}

	topic := &Topic{
		TopicID: uuid.Must(uuid.NewV4()).String(),
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
		count, err := topicsCountByCategory(ctx, tx, category.CategoryID)
		if err != nil {
			return err
		}
		category.TopicsCount = count + 1
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
func (user *User) UpdateTopic(ctx context.Context, id, title, body, categoryID string) (*Topic, error) {
	title, body = strings.TrimSpace(title), strings.TrimSpace(body)
	if title != "" && len(title) < minTitleSize {
		return nil, session.BadDataError(ctx)
	}

	var topic *Topic
	var prevCategoryID string
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
		if categoryID != "" && topic.CategoryID != categoryID {
			prevCategoryID = topic.CategoryID
			category, err := findCategory(ctx, tx, categoryID)
			if err != nil {
				return err
			} else if category == nil {
				return session.NotFoundError(ctx)
			}
			topic.CategoryID = category.CategoryID
			topic.Category = category
		}
		topic.User = user
		return tx.Update(topic)
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return nil, session.TransactionError(ctx, err)
	}
	if prevCategoryID != "" {
		go ElevateCategory(ctx, prevCategoryID)
		go ElevateCategory(ctx, topic.CategoryID)
	}
	return topic, nil
}

//ReadTopic read a topic by ID
func ReadTopic(ctx context.Context, id string) (*Topic, error) {
	topic := &Topic{TopicID: id}
	if err := session.Database(ctx).Model(topic).Relation("User").Relation("Category").WherePK().Select(); err == pg.ErrNoRows {
		return nil, session.NotFoundError(ctx)
	} else if err != nil {
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
	if err := session.Database(ctx).Model(&topics).Relation("User").Relation("Category").Where("topic.created_at<?", offset).Order("topic.created_at DESC").Limit(50).Select(); err != nil {
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
	if err := session.Database(ctx).Model(&topics).Relation("Category").Where("topic.user_id=? AND topic.created_at<?", user.UserID, offset).Order("topic.created_at DESC").Limit(50).Select(); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	for _, topic := range topics {
		topic.User = user
	}
	return topics, nil
}

// ReadTopics read topics by CategoryID order by created_at DESC
func (category *Category) ReadTopics(ctx context.Context, offset time.Time) ([]*Topic, error) {
	if offset.IsZero() {
		offset = time.Now()
	}
	var topics []*Topic
	if err := session.Database(ctx).Model(&topics).Relation("User").Where("topic.category_id=? AND topic.created_at<?", category.CategoryID, offset).Order("topic.created_at DESC").Limit(50).Select(); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	for _, topic := range topics {
		topic.Category = category
	}
	return topics, nil
}

func (category *Category) lastTopic(ctx context.Context, tx *pg.Tx) (*Topic, error) {
	var topic Topic
	if err := tx.Model(&topic).Where("category_id=?", category.CategoryID).Order("created_at DESC").Limit(1).Select(); err != nil {
		return nil, err
	}
	return &topic, nil
}

func topicsCountByCategory(ctx context.Context, tx *pg.Tx, id string) (int, error) {
	return tx.Model(&Topic{}).Where("category_id=?", id).Count()
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
