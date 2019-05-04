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
	title, body := strings.TrimSpace(p.Title), strings.TrimSpace(p.Body)
	if title != "" && len(title) < minTitleSize {
		return nil, session.BadDataError(ctx)
	}

	var topic *Model
	var prevCategoryID string
	err := t.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		topic, err = findTopic(ctx, tx, id)
		if err != nil {
			return err
		} else if topic == nil {
			return nil
		} else if topic.UserID != user.UserID && !user.IsAdmin() {
			return session.AuthorizationError(ctx)
		}
		if title != "" {
			topic.Title = title
		}
		topic.Body = body
		if p.CategoryID != "" && topic.CategoryID != p.CategoryID {
			prevCategoryID = topic.CategoryID
			// todo: use public category function
			category, err := t.categoryStore.Find(ctx, tx, p.CategoryID)
			if err != nil {
				return err
			} else if category == nil {
				return session.BadDataError(ctx)
			}
			topic.CategoryID = category.CategoryID
			topic.Category = category
		}
		cols, params := durable.PrepareColumnsWithValues([]string{"title", "body", "category_id"})
		vals := []interface{}{topic.Title, topic.Body, topic.CategoryID}
		_, err = tx.ExecContext(ctx, fmt.Sprintf("UPDATE topics SET (%s)=(%s) WHERE topic_id='%s'", cols, params, topic.TopicID), vals...)
		return err
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return nil, session.TransactionError(ctx, err)
	}
	if topic == nil {
		return nil, session.NotFoundError(ctx)
	}
	if prevCategoryID != "" {
		go t.dispersalCategory(ctx, prevCategoryID)
		go t.dispersalCategory(ctx, topic.CategoryID)
	}
	topic.User = user
	return topic, nil
}

func (t *Topic) GetByID(ctx context.Context, id string) (*Model, error) {
	var topic *Model
	err := t.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		topic, err = findTopic(ctx, tx, id)
		if err != nil {
			return err
		}
		if topic == nil {
			subs := strings.Split(id, "-")
			if len(subs) < 1 || len(subs[0]) <= 5 {
				return nil
			}
			id = subs[0]
			topic, err = findTopicByShortID(ctx, tx, id)
			if topic == nil || err != nil {
				return err
			}
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
	if offset.IsZero() {
		offset = time.Now()
	}

	var topics []*Model
	err := t.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		set, err := t.categoryStore.GetCategorySet(ctx, tx)
		if err != nil {
			return err
		}

		query := fmt.Sprintf("SELECT %s FROM topics WHERE created_at<$1 ORDER BY created_at DESC LIMIT $2", strings.Join(topicColumns, ","))
		rows, err := tx.QueryContext(ctx, query, offset, LIMIT)
		if err != nil {
			return err
		}
		defer rows.Close()

		userIds := []string{}
		for rows.Next() {
			topic, err := topicFromRows(rows)
			if err != nil {
				return err
			}
			userIds = append(userIds, topic.UserID)
			topic.Category = set[topic.CategoryID]
			topics = append(topics, topic)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		userSet, err := t.userStore.GetUserSet(ctx, tx, userIds)
		if err != nil {
			return err
		}
		for i, topic := range topics {
			topics[i].User = userSet[topic.UserID]
		}
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topics, nil
}

func (t *Topic) GetByUserID(ctx context.Context, user *user.Model, offset time.Time) ([]*Model, error) {
	if offset.IsZero() {
		offset = time.Now()
	}

	var topics []*Model
	err := t.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		set, err := t.categoryStore.GetCategorySet(ctx, tx)
		if err != nil {
			return err
		}
		query := fmt.Sprintf("SELECT %s FROM topics WHERE user_id=$1 AND created_at<$2 ORDER BY created_at DESC LIMIT $3", strings.Join(topicColumns, ","))
		rows, err := tx.QueryContext(ctx, query, user.UserID, offset, LIMIT)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			topic, err := topicFromRows(rows)
			if err != nil {
				return err
			}
			topic.User = user
			topic.Category = set[topic.CategoryID]
			topics = append(topics, topic)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topics, nil
}

func (t *Topic) GetByCategoryID(ctx context.Context, cat *category.Model, offset time.Time) ([]*Model, error) {
	if offset.IsZero() {
		offset = time.Now()
	}

	var topics []*Model
	err := t.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		query := fmt.Sprintf("SELECT %s FROM topics WHERE category_id=$1 AND created_at<$2 ORDER BY created_at DESC LIMIT $3", strings.Join(topicColumns, ","))
		rows, err := tx.QueryContext(ctx, query, cat.CategoryID, offset, LIMIT)
		if err != nil {
			return err
		}
		defer rows.Close()

		userIds := []string{}
		for rows.Next() {
			topic, err := topicFromRows(rows)
			if err != nil {
				return err
			}
			userIds = append(userIds, topic.UserID)
			topic.Category = cat
			topics = append(topics, topic)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		userSet, err := t.userStore.GetUserSet(ctx, tx, userIds)
		if err != nil {
			return err
		}
		for i, topic := range topics {
			topics[i].User = userSet[topic.UserID]
		}
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topics, nil
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
