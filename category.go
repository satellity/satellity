package models

import (
	"context"
	"database/sql"
	"fmt"
	"godiscourse/internal/durable"
	"godiscourse/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

// topics_count should use pg int64
const categoriesDDL = `
CREATE TABLE IF NOT EXISTS categories (
	category_id           VARCHAR(36) PRIMARY KEY,
	name                  VARCHAR(36) NOT NULL,
	alias                 VARCHAR(128) NOT NULL,
	description           VARCHAR(512) NOT NULL,
	topics_count          INTEGER NOT NULL DEFAULT 0,
	last_topic_id         VARCHAR(36),
	position              INTEGER NOT NULL DEFAULT 0,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX ON categories (position);
`

// Category is used to categorize topics.
type Category struct {
	CategoryID  string
	Name        string
	Alias       string
	Description string
	TopicsCount int64
	LastTopicID sql.NullString
	Position    int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var categoryColumns = []string{"category_id", "name", "alias", "description", "topics_count", "last_topic_id", "position", "created_at", "updated_at"}

func (c *Category) values() []interface{} {
	return []interface{}{c.CategoryID, c.Name, c.Alias, c.Description, c.TopicsCount, c.LastTopicID, c.Position, c.CreatedAt, c.UpdatedAt}
}

// CreateCategory create a new category.
func CreateCategory(mctx *Context, name, alias, description string, position int64) (*Category, error) {
	ctx := mctx.context
	alias, name = strings.TrimSpace(alias), strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if len(name) < 1 {
		return nil, session.BadDataError(ctx)
	}
	if alias == "" {
		alias = name
	}

	t := time.Now()
	category := &Category{
		CategoryID:  uuid.Must(uuid.NewV4()).String(),
		Name:        name,
		Alias:       alias,
		Description: description,
		TopicsCount: 0,
		LastTopicID: sql.NullString{String: "", Valid: false},
		Position:    position,
		CreatedAt:   t,
		UpdatedAt:   t,
	}

	cols, params := durable.PrepareColumnsWithValues(categoryColumns)
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		if position == 0 {
			count, err := categoryCount(ctx, tx)
			if err != nil {
				return err
			}
			category.Position = count
		}
		_, err := tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO categories(%s) VALUES (%s)", cols, params), category.values()...)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return category, nil
}

// UpdateCategory update a category's attributes
func UpdateCategory(mctx *Context, id, name, alias, description string, position int64) (*Category, error) {
	ctx := mctx.context
	alias, name = strings.TrimSpace(alias), strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if len(alias) < 1 && len(name) < 1 {
		return nil, session.BadDataError(ctx)
	}

	var category *Category
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		category, err = findCategory(ctx, tx, id)
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
		category.Position = position
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

// ReadCategory read a category by ID (uuid).
func ReadCategory(mctx *Context, id string) (*Category, error) {
	ctx := mctx.context
	var category *Category
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		category, err = findCategory(ctx, tx, id)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return category, nil
}

// ReadAllCategories read categories order by position
func ReadAllCategories(mctx *Context) ([]*Category, error) {
	ctx := mctx.context
	var categories []*Category
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		categories, err = readCategories(ctx, tx)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return categories, nil
}

func readCategorySet(ctx context.Context, tx *sql.Tx) (map[string]*Category, error) {
	categories, err := readCategories(ctx, tx)
	if err != nil {
		return nil, err
	}
	set := make(map[string]*Category, 0)
	for _, c := range categories {
		set[c.CategoryID] = c
	}
	return set, nil
}

func readCategories(ctx context.Context, tx *sql.Tx) ([]*Category, error) {
	rows, err := tx.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM categories ORDER BY position LIMIT 500", strings.Join(categoryColumns, ",")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		category, err := categoryFromRows(rows)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

// dispersalCategory update category's info, e.g.: LastTopicID, TopicsCount
func dispersalCategory(mctx *Context, id string) (*Category, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}
	ctx := context.Background()
	var category *Category
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		category, err = findCategory(ctx, tx, id)
		if err != nil {
			return err
		} else if category == nil {
			return session.NotFoundError(ctx)
		}
		topic, err := category.lastTopic(ctx, tx)
		if err != nil {
			return err
		}
		var lastTopicID = sql.NullString{String: "", Valid: false}
		if topic != nil {
			lastTopicID = sql.NullString{String: topic.TopicID, Valid: true}
		}
		if category.LastTopicID.String != lastTopicID.String {
			category.LastTopicID = lastTopicID
		}
		category.TopicsCount = 0
		if category.LastTopicID.Valid {
			count, err := topicsCountByCategory(ctx, tx, category.CategoryID)
			if err != nil {
				return err
			}
			category.TopicsCount = count
		}
		category.UpdatedAt = time.Now()
		cols, params := durable.PrepareColumnsWithValues([]string{"last_topic_id", "topics_count", "updated_at"})
		vals := []interface{}{category.LastTopicID, category.TopicsCount, category.UpdatedAt}
		_, err = tx.ExecContext(ctx, fmt.Sprintf("UPDATE categories SET (%s)=(%s) WHERE category_id='%s'", cols, params, category.CategoryID), vals...)
		return err
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return nil, session.TransactionError(ctx, err)
	}
	return category, nil
}

func findCategory(ctx context.Context, tx *sql.Tx, id string) (*Category, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}

	row := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM categories WHERE category_id=$1", strings.Join(categoryColumns, ",")), id)
	c, err := categoryFromRows(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}

func categoryCount(ctx context.Context, tx *sql.Tx) (int64, error) {
	var count int64
	row := tx.QueryRowContext(ctx, "SELECT count(*) FROM categories")
	err := row.Scan(&count)
	if err != nil {
		return 0, session.TransactionError(ctx, err)
	}
	return count, nil
}

func categoryFromRows(row durable.Row) (*Category, error) {
	var c Category
	err := row.Scan(&c.CategoryID, &c.Name, &c.Alias, &c.Description, &c.TopicsCount, &c.LastTopicID, &c.Position, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}
