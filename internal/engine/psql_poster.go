package engine

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"godiscourse/internal/durable"
	"godiscourse/internal/model"
	"godiscourse/internal/session"

	"github.com/gofrs/uuid"
)

type Psql struct {
	db *durable.Database
}

func NewPsql(db *durable.Database) *Psql {
	return &Psql{db: db}
}

func (p *Psql) GetCategoryByID(ctx context.Context, id string) (*model.Category, error) {
	var category *model.Category
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		category, err = model.FindCategory(ctx, tx, id)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return category, nil
}

func (p *Psql) CreateTopic(ctx context.Context, userID string, t *model.TopicInfo) (*model.Topic, error) {
	title, body := strings.TrimSpace(t.Title), strings.TrimSpace(t.Body)
	if len(title) < minTitleSize {
		return nil, session.BadDataError(ctx)
	}

	now := time.Now()
	topic := &model.Topic{
		TopicID:   uuid.Must(uuid.NewV4()).String(),
		Title:     title,
		Body:      body,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	var err error
	topic.ShortID, err = model.GenerateShortID("topics", now)
	if err != nil {
		return nil, session.ServerError(ctx, err)
	}

	err = p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		category, err := p.GetCategoryByID(ctx, t.CategoryID)
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
		cols, params := durable.PrepareColumnsWithValues(model.TopicColumns)
		_, err = tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO topics(%s) VALUES (%s)", cols, params), topic.Values()...)
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

func (p *Psql) UpdateTopic(ctx context.Context, id string, t *model.TopicInfo) (*model.Topic, error) {
	title, body := strings.TrimSpace(t.Title), strings.TrimSpace(t.Body)
	if title != "" && len(title) < minTitleSize {
		return nil, session.BadDataError(ctx)
	}

	var topic *model.Topic
	var prevCategoryID string
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		topic, err = model.FindTopic(ctx, tx, id)
		if err != nil {
			return err
		} else if topic == nil {
			return nil
		}
		// todo: move to level up
		// } else if topic.UserID != user.UserID && !user.IsAdmin() {
		// 	return session.AuthorizationError(ctx)
		// }
		if title != "" {
			topic.Title = title
		}
		topic.Body = body
		if t.CategoryID != "" && topic.CategoryID != t.CategoryID {
			prevCategoryID = topic.CategoryID
			category, err := model.FindCategory(ctx, tx, t.CategoryID)
			if err != nil {
				return err
			} else if category == nil {
				return session.BadDataError(ctx)
			}
			topic.CategoryID = category.CategoryID
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
		// go t.dispersalCategory(ctx, prevCategoryID)
		// go t.dispersalCategory(ctx, topic.CategoryID)
	}
	return topic, nil
}

// todo: rewrite with join
func (p *Psql) GetTopicByID(ctx context.Context, id string) (*model.Topic, error) {
	var topic *model.Topic
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		topic, err = model.FindTopic(ctx, tx, id)
		if err != nil {
			return err
		}
		if topic == nil {
			subs := strings.Split(id, "-")
			if len(subs) < 1 || len(subs[0]) <= 5 {
				return nil
			}
			id = subs[0]
			topic, err = model.FindTopicByShortID(ctx, tx, id)
			if topic == nil || err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topic, nil
}

func (p *Psql) GetTopicByUserID(ctx context.Context, userID string, offset time.Time) ([]*model.Topic, error) {
	if offset.IsZero() {
		offset = time.Now()
	}

	var topics []*model.Topic
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		// todo: join query
		query := fmt.Sprintf("SELECT %s FROM topics WHERE user_id=$1 AND created_at<$2 ORDER BY created_at DESC LIMIT $3", strings.Join(model.TopicColumns, ","))
		rows, err := tx.QueryContext(ctx, query, userID, offset, LIMIT)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			topic, err := model.TopicFromRows(rows)
			if err != nil {
				return err
			}
			topics = append(topics, topic)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topics, nil
}

func (p *Psql) GetByCategoryID(ctx context.Context, categoryID string, offset time.Time) ([]*model.Topic, error) {
	if offset.IsZero() {
		offset = time.Now()
	}

	var topics []*model.Topic
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		// todo: join query
		query := fmt.Sprintf("SELECT %s FROM topics WHERE category_id=$1 AND created_at<$2 ORDER BY created_at DESC LIMIT $3", strings.Join(model.TopicColumns, ","))
		rows, err := tx.QueryContext(ctx, query, categoryID, offset, LIMIT)
		if err != nil {
			return err
		}
		defer rows.Close()
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topics, nil
}

func (p *Psql) GetTopicsByOffset(ctx context.Context, offset time.Time) ([]*model.Topic, error) {
	if offset.IsZero() {
		offset = time.Now()
	}

	var topics []*model.Topic
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		// todo: join query
		query := fmt.Sprintf("SELECT %s FROM topics WHERE created_at<$1 ORDER BY created_at DESC LIMIT $2", strings.Join(model.TopicColumns, ","))
		rows, err := tx.QueryContext(ctx, query, offset, LIMIT)
		if err != nil {
			return err
		}
		defer rows.Close()
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topics, nil
}
