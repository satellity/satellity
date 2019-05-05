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

	User     *user.Model
	Category *category.Model
}

type Params struct {
	Title      string
	Body       string
	CategoryID string
}

type TopicDatastore interface {
	Create(ctx context.Context, uid string, p *Params) (*Model, error)
	Update(ctx context.Context, user *user.Model, id string, p *Params) (*Model, error)
	GetByID(ctx context.Context, id string) (*Model, error)
	GetByShortID(ctx context.Context, id string) (*Model, error)
	GetByOffset(ctx context.Context, offset time.Time) ([]*Model, error) // equal ReadTopics
	GetByUserID(ctx context.Context, user *user.Model, offset time.Time) ([]*Model, error)
	GetByCategoryID(ctx context.Context, cat *category.Model, offset time.Time) ([]*Model, error)
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
	return &Model{}, nil
}

func (t *Topic) Update(ctx context.Context, user *user.Model, id string, p *Params) (*Model, error) {
	return &Model{}, nil
	}

func (t *Topic) GetByID(ctx context.Context, id string) (*Model, error) {
	return &Model{}, nil
			}

func (t *Topic) GetByShortID(ctx context.Context, id string) (*Model, error) {
	subs := strings.Split(id, "-")
	if len(subs) < 1 || len(subs[0]) <= 5 {
		return nil, nil
	}
	id = subs[0]

	var topic *Model
	err := t.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		topic, err = findTopicByShortID(ctx, tx, id)
		if topic == nil || err != nil {
			return err
		}
		user, err := t.userStore.GetByID(ctx, topic.UserID)
		if err != nil {
			return err
		}
		category, err := t.categoryStore.Find(ctx, tx, topic.CategoryID)
		if err != nil {
			return err
		}
		topic.User = user
		topic.Category = category
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topic, nil
}

func (t *Topic) GetByOffset(ctx context.Context, offset time.Time) ([]*Model, error) {
	return []*Model{}, nil
		}

func (t *Topic) GetByUserID(ctx context.Context, user *user.Model, offset time.Time) ([]*Model, error) {
	return []*Model{}, nil
		}

func (t *Topic) GetByCategoryID(ctx context.Context, cat *category.Model, offset time.Time) ([]*Model, error) {
	return []*Model{}, nil
}

func LastTopic(ctx context.Context, cid string, tx *sql.Tx) (*Model, error) {
	row := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM topics WHERE category_id=$1 LIMIT 1", strings.Join(topicColumns, ",")), cid)
	t, err := topicFromRows(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return t, err
}

// dispersalCategory update category's info, e.g.: LastTopicID, TopicsCount
func (t *Topic) dispersalCategory(ctx context.Context, id string) (*category.Model, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}
	var result *category.Model
	err := t.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		result, err = t.categoryStore.Find(ctx, tx, id)
		if err != nil {
			return err
		} else if result == nil {
			return session.NotFoundError(ctx)
		}
		topic, err := LastTopic(ctx, result.CategoryID, tx)
		if err != nil {
			return err
		}
		var lastTopicID = sql.NullString{String: "", Valid: false}
		if topic != nil {
			lastTopicID = sql.NullString{String: topic.TopicID, Valid: true}
		}
		if result.LastTopicID.String != lastTopicID.String {
			result.LastTopicID = lastTopicID
		}
		result.TopicsCount = 0
		if result.LastTopicID.Valid {
			count, err := topicsCountByCategory(ctx, tx, result.CategoryID)
			if err != nil {
				return err
			}
			result.TopicsCount = count
		}
		result.UpdatedAt = time.Now()
		cols, params := durable.PrepareColumnsWithValues([]string{"last_topic_id", "topics_count", "updated_at"})
		vals := []interface{}{result.LastTopicID, result.TopicsCount, result.UpdatedAt}
		_, err = tx.ExecContext(ctx, fmt.Sprintf("UPDATE categories SET (%s)=(%s) WHERE category_id='%s'", cols, params, result.CategoryID), vals...)
		return err
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return nil, session.TransactionError(ctx, err)
	}
	return result, nil
}
