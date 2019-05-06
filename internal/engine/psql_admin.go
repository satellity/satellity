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

func (p *Psql) GetUsersByOffset(ctx context.Context, offset time.Time) ([]*models.User, error) {
	if offset.IsZero() {
		offset = time.Now()
	}
	rows, err := p.db.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM users WHERE created_at<$1 ORDER BY created_at DESC LIMIT 100", strings.Join(models.UserColumns, ",")), offset)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user, err := models.UserFromRows(rows)
		if err != nil {
			return nil, session.TransactionError(ctx, err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return users, nil
}

func (p *Psql) CreateCategory(ctx context.Context, c *models.CategoryInfo) (*models.Category, error) {
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
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
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

func (p *Psql) UpdateCategory(ctx context.Context, id string, c *models.CategoryInfo) (*models.Category, error) {
	alias, name := strings.TrimSpace(c.Alias), strings.TrimSpace(c.Name)
	description := strings.TrimSpace(c.Description)
	if len(alias) < 1 && len(name) < 1 {
		return nil, session.BadDataError(ctx)
	}

	var category *models.Category
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		category, err := models.FindCategory(ctx, tx, id)
		if err != nil || category == nil {
			return err
		}
		if len(name) > 0 {
			category.Name = name
		}
		if len(alias) > 0 {
			category.Alias = alias
		}
		if len(description) > 0 {
			category.Description = description
		}
		category.Position = c.Position
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
