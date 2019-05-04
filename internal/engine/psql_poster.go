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
			// todo: use public category function
			category, err := model.FindCategory(ctx, tx, t.CategoryID)
			if err != nil {
				return err
			} else if category == nil {
				return session.BadDataError(ctx)
			}
			topic.CategoryID = category.CategoryID
			// topic.Category = category
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
	// topic.User = user
	return topic, nil
}
