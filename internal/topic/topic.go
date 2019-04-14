package topic

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"godiscourse/internal/category"
	"godiscourse/internal/durable"
	"godiscourse/internal/session"
	"godiscourse/internal/user"

	"github.com/gofrs/uuid"
)

// Topic related CONST
const (
	minTitleSize = 3
	LIMIT        = 50
)

type Model struct {
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

	User     *user.Data
	Category *category.Model
}

type Params struct {
	Title      string
	Body       string
	CategoryID string
}

type TopicDatastore interface {
	Create(ctx context.Context, uid string, p *Params) (*Model, error)
	Update(ctx context.Context, uid string, id string, p *Params) (*Model, error)
	GetByID(ctx context.Context, id string) (*Model, error)
	GetByShortID(ctx context.Context, id string) (*Model, error)
	GetByOffset(ctx context.Context, offset time.Time) ([]*Model, error)
	GetByUserID(ctx context.Context, uid string, offset time.Time) ([]*Model, error)
	GetByCategoryID(ctx context.Context, cid string, offset time.Time) ([]*Model, error)
}

type Topic struct {
	db            *durable.Database
	userStore     user.UserDatastore
	categoryStore category.CategoryDatastore
}

func New(db *durable.Database, u user.UserDatastore, c category.CategoryDatastore) *Topic {
	return &Topic{
		db:            db,
		userStore:     u,
		categoryStore: c,
	}
}

func (t *Topic) Create(ctx context.Context, uid string, p *Params) (*Model, error) {
	title, body := strings.TrimSpace(p.Title), strings.TrimSpace(p.Body)
	if len(title) < minTitleSize {
		return nil, session.BadDataError(ctx)
	}

	now := time.Now()
	topic := &Model{
		TopicID:   uuid.Must(uuid.NewV4()).String(),
		Title:     title,
		Body:      body,
		UserID:    uid,
		CreatedAt: now,
		UpdatedAt: now,
	}

	var err error
	topic.ShortID, err = generateShortID("topics", now)
	if err != nil {
		return nil, session.ServerError(ctx, err)
	}

	err = t.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		category, err := t.categoryStore.GetByID(ctx, p.CategoryID)
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
		category.TopicsCount, category.UpdatedAt = count+1, time.Now()
		cols, params := durable.PrepareColumnsWithValues(topicColumns)
		_, err = tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO topics(%s) VALUES (%s)", cols, params), topic.values()...)
		if err != nil {
			return err
		}
		ccols, cparams := durable.PrepareColumnsWithValues([]string{"last_topic_id", "topics_count", "updated_at"})
		cvals := []interface{}{category.LastTopicID, category.TopicsCount, category.UpdatedAt}
		_, err = tx.ExecContext(ctx, fmt.Sprintf("UPDATE categories SET (%s)=(%s) WHERE category_id='%s'", ccols, cparams, category.CategoryID), cvals...)
		if err != nil {
			return err
		}
		// _, err = upsertStatistic(ctx, tx, "topics")
		return err
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return nil, session.TransactionError(ctx, err)
	}
	return topic, nil
}
