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

const categories_DDL = `
CREATE TABLE IF NOT EXISTS categories (
	category_id           VARCHAR(36) PRIMARY KEY,
	name                  VARCHAR(36) NOT NULL,
	description           VARCHAR(512) NOT NULL,
	topics_count          INTEGER NOT NULL DEFAULT 0,
	last_topic_id         VARCHAR(36),
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX ON categories (name);
`

var categoryCols = []string{"category_id", "name", "description", "topics_count", "last_topic_id", "created_at", "updated_at"}

type Category struct {
	CategoryId  string         `sql:"category_id,pk"`
	Name        string         `sql:"name"`
	Description string         `sql:"description"`
	TopicsCount int            `sql:"topics_count"`
	LastTopicId sql.NullString `sql:"last_topic_id"`
	CreatedAt   time.Time      `sql:"created_at"`
	UpdatedAt   time.Time      `sql:"updated_at"`
}

func CreateCategory(ctx context.Context, name, description string) (*Category, error) {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)

	t := time.Now()
	category := &Category{
		CategoryId:  uuid.NewV4().String(),
		Name:        name,
		Description: description,
		LastTopicId: sql.NullString{"", false},
		CreatedAt:   t,
		UpdatedAt:   t,
	}

	if err := session.Database(ctx).Insert(category); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return category, nil
}

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

func findCategory(ctx context.Context, tx *pg.Tx, id string) (*Category, error) {
	category := &Category{CategoryId: id}
	if err := tx.Model(category).Column(categoryCols...).WherePK().Select(); err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return category, nil
}
