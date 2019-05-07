package engine

import (
	"context"
	"database/sql"
	"fmt"
	"godiscourse/internal/durable"
	"godiscourse/internal/models"
	"godiscourse/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

func (s *Store) CreateCategory(ctx context.Context, c *models.CategoryRequest) (*models.Category, error) {
	alias, name := strings.TrimSpace(c.Alias), strings.TrimSpace(c.Name)
	description := strings.TrimSpace(c.Description)
	if len(name) < 1 {
		return nil, session.BadDataError(ctx)
	}
	if alias == "" {
		alias = c.Name
	}

	t := time.Now()
	category := &models.Category{
		CategoryID:  uuid.Must(uuid.NewV4()).String(),
		Name:        name,
		Alias:       alias,
		Description: description,
		TopicsCount: 0,
		LastTopicID: sql.NullString{String: "", Valid: false},
		Position:    c.Position,
		CreatedAt:   t,
		UpdatedAt:   t,
	}

	cols, params := durable.PrepareColumnsWithValues(models.CategoryColumns)
	err := s.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		if c.Position == 0 {
			count, err := categoryCount(ctx, tx)
			if err != nil {
				return err
			}
			category.Position = count
		}
		_, err := tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO categories(%s) VALUES (%s)", cols, params), category.Values()...)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return category, nil
}

func (s *Store) UpdateCategory(ctx context.Context, id string, c *models.CategoryRequest) (*models.Category, error) {
	c.Alias = strings.TrimSpace(c.Alias)
	c.Name = strings.TrimSpace(c.Name)
	c.Description = strings.TrimSpace(c.Description)
	if len(c.Alias) < 1 && len(c.Name) < 1 {
		return nil, session.BadDataError(ctx)
	}

	var category *models.Category
	err := s.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		category, err = models.FindCategory(ctx, tx, id)
		if err != nil || category == nil {
			return err
		}
		if len(c.Name) > 0 {
			category.Name = c.Name
		}
		if len(c.Alias) > 0 {
			category.Alias = c.Alias
		}
		if len(c.Description) > 0 {
			category.Description = c.Description
		}
		if c.Position > 0 {
			category.Position = c.Position
		}
		category.UpdatedAt = time.Now()
		cols, params := durable.PrepareColumnsWithValues([]string{"name", "alias", "description", "position", "updated_at"})
		vals := []interface{}{category.Name, category.Alias, category.Description, category.Position, category.UpdatedAt}
		_, err = tx.ExecContext(ctx, fmt.Sprintf("UPDATE categories SET (%s)=(%s) WHERE category_id='%s'", cols, params, category.CategoryID), vals...)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	if category == nil {
		return nil, session.NotFoundError(ctx)
	}
	return category, nil
}

func (s *Store) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	var categories []*models.Category
	err := s.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		categories, err = models.ReadCategories(ctx, tx)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return categories, nil
}

// transmitToCategory update category's info, e.g.: LastTopicID, TopicsCount
func (s *Store) transmitToCategory(ctx context.Context, id string) (*models.Category, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}
	var result *models.Category
	err := s.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		result, err = models.FindCategory(ctx, tx, id)
		if err != nil {
			return err
		} else if result == nil {
			return session.NotFoundError(ctx)
		}
		topic, err := models.LastTopic(ctx, result.CategoryID, tx)
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
