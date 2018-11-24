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

const categoriesDDL = `
CREATE TABLE IF NOT EXISTS categories (
	category_id           VARCHAR(36) PRIMARY KEY,
	name                  VARCHAR(36) NOT NULL,
	description           VARCHAR(512) NOT NULL,
	topics_count          INTEGER NOT NULL DEFAULT 0,
	last_topic_id         VARCHAR(36),
	position              INTEGER NOT NULL DEFAULT 0,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX ON categories (position);
`

var categoryCols = []string{"category_id", "name", "description", "topics_count", "last_topic_id", "position", "created_at", "updated_at"}

// Category is used to categorize topics.
type Category struct {
	CategoryID  string         `sql:"category_id,pk"`
	Name        string         `sql:"name"`
	Description string         `sql:"description"`
	TopicsCount int            `sql:"topics_count,notnull"`
	LastTopicID sql.NullString `sql:"last_topic_id"`
	Position    int            `sql:"position,notnull"`
	CreatedAt   time.Time      `sql:"created_at"`
	UpdatedAt   time.Time      `sql:"updated_at"`
}

// CreateCategory create a new category.
func CreateCategory(ctx context.Context, name, description string) (*Category, error) {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if len(name) < 1 {
		return nil, session.BadDataError(ctx)
	}

	t := time.Now()
	count, err := categoryCount(ctx)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	category := &Category{
		CategoryID:  uuid.NewV4().String(),
		Name:        name,
		Description: description,
		LastTopicID: sql.NullString{String: "", Valid: false},
		Position:    count,
		CreatedAt:   t,
		UpdatedAt:   t,
	}

	if err := session.Database(ctx).Insert(category); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return category, nil
}

// UpdateCategory update a category
func UpdateCategory(ctx context.Context, id, name, description string) (*Category, error) {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if len(name) < 1 && len(description) < 1 {
		return nil, session.BadDataError(ctx)
	}
	var category *Category
	err := session.Database(ctx).RunInTransaction(func(tx *pg.Tx) error {
		var err error
		category, err = findCategory(ctx, tx, id)
		if err != nil {
			return err
		}
		if len(name) > 0 {
			category.Name = name
		}
		if len(description) > 0 {
			category.Description = description
		}
		category.UpdatedAt = time.Now()
		return tx.Update(category)
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return category, nil
}

// ReadCategory read a category by ID (uuid).
func ReadCategory(ctx context.Context, id string) (*Category, error) {
	var category *Category
	err := session.Database(ctx).RunInTransaction(func(tx *pg.Tx) error {
		var err error
		category, err = findCategory(ctx, tx, id)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return category, nil
}

// ReadCategories read categories
func ReadCategories(ctx context.Context) ([]*Category, error) {
	var categories []*Category
	if _, err := session.Database(ctx).Query(&categories, "SELECT * FROM categories ORDER BY position LIMIT 100"); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return categories, nil
}

func findCategory(ctx context.Context, tx *pg.Tx, id string) (*Category, error) {
	category := &Category{CategoryID: id}
	if err := tx.Model(category).Column(categoryCols...).WherePK().Select(); err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return category, nil
}

func categoryCount(ctx context.Context) (int, error) {
	return session.Database(ctx).Model(&Category{}).Count()
}
