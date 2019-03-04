package models

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/godiscourse/godiscourse/api/session"
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

var categoryCols = []string{"category_id", "name", "alias", "description", "topics_count", "last_topic_id", "position", "created_at", "updated_at"}

func (c *Category) values() []interface{} {
	return []interface{}{c.CategoryID, c.Name, c.Alias, c.Description, c.TopicsCount, c.LastTopicID, c.Position, c.CreatedAt, c.UpdatedAt}
}

// CreateCategory create a new category.
func CreateCategory(ctx context.Context, name, alias, description string, position int64) (*Category, error) {
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
	if position == 0 {
		count, err := categoryCount(ctx)
		if err != nil {
			return nil, err
		}
		category.Position = count
	}

	cols, params := prepareColumnsWithValues(categoryCols)
	_, err := session.Database(ctx).ExecContext(ctx, fmt.Sprintf("INSERT INTO categories(%s) VALUES (%s)", cols, params), category.values()...)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return category, nil
}

// UpdateCategory update a category's attributes
func UpdateCategory(ctx context.Context, id, name, alias, description string, position int64) (*Category, error) {
	alias, name = strings.TrimSpace(alias), strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if len(alias) < 1 && len(name) < 1 {
		return nil, session.BadDataError(ctx)
	}

	var category *Category
	err := runInTransaction(ctx, func(tx *sql.Tx) error {
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
		cols, params := prepareColumnsWithValues([]string{"name", "alias", "description", "position", "updated_at"})
		vals := []interface{}{category.Name, category.Alias, category.Description, category.Position, category.UpdatedAt}
		_, err = tx.ExecContext(ctx, fmt.Sprintf("UPDATE categories SET (%s)=(%s) WHERE category_id='%s'", cols, params, category.CategoryID), vals...)
		return err
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return nil, session.TransactionError(ctx, err)
	}
	if category == nil {
		return nil, session.NotFoundError(ctx)
	}
	return category, nil
}

// ReadCategory read a category by ID (uuid).
func ReadCategory(ctx context.Context, id string) (*Category, error) {
	var category *Category
	err := runInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		category, err = findCategory(ctx, tx, id)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return category, nil
}

// ReadCategories read categories order by position
func ReadCategories(ctx context.Context) ([]*Category, error) {
	rows, err := session.Database(ctx).QueryContext(ctx, fmt.Sprintf("SELECT %s FROM categories ORDER BY position LIMIT 100", strings.Join(categoryCols, ",")))
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		category, err := categoryFromRows(rows)
		if err != nil {
			return nil, session.TransactionError(ctx, err)
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return categories, nil
}

func readCategorySet(ctx context.Context) (map[string]*Category, error) {
	categories, err := ReadCategories(ctx)
	if err != nil {
		return nil, err
	}
	set := make(map[string]*Category, 0)
	for _, c := range categories {
		set[c.CategoryID] = c
	}
	return set, nil
}

// ElevateCategory update category's info, e.g.: LastTopicID, TopicsCount
func ElevateCategory(ctx context.Context, id string) (*Category, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}
	var category *Category
	err := runInTransaction(ctx, func(tx *sql.Tx) error {
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
		var lastTopicId = sql.NullString{String: "", Valid: false}
		if topic != nil {
			lastTopicId = sql.NullString{String: topic.TopicID, Valid: true}
		}
		if category.LastTopicID.String != lastTopicId.String {
			category.LastTopicID = lastTopicId
		}
		category.TopicsCount = 0
		if category.LastTopicID.Valid {
			count, err := topicsCountByCategory(ctx, tx, category.CategoryID)
			if err != nil {
				return err
			}
			category.UpdatedAt = time.Now()
			cols, params := prepareColumnsWithValues([]string{"last_topic_id", "topics_count", "updated_at"})
			vals := []interface{}{category.LastTopicID, category.TopicsCount, category.UpdatedAt}
			_, err = tx.ExecContext(ctx, fmt.Sprintf("UPDATE categories SET (%s)=(%s) WHERE category_id='%s'", cols, params, category.CategoryID), vals...)
			category.TopicsCount = count
		}
		return nil
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

	rows, err := tx.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM categories WHERE category_id=$1", strings.Join(categoryCols, ",")), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, nil
	}
	c, err := categoryFromRows(rows)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func categoryCount(ctx context.Context) (int64, error) {
	var count int64
	err := session.Database(ctx).QueryRowContext(ctx, "SELECT count(*) FROM categories").Scan(&count)
	if err != nil {
		return 0, session.TransactionError(ctx, err)
	}
	return count, nil
}

func categoryFromRows(rows *sql.Rows) (*Category, error) {
	var c Category
	err := rows.Scan(&c.CategoryID, &c.Name, &c.Alias, &c.Description, &c.TopicsCount, &c.LastTopicID, &c.Position, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func prepareColumnsWithValues(columns []string) (string, string) {
	if len(columns) < 1 {
		return "", ""
	}
	cols, params := bytes.Buffer{}, bytes.Buffer{}
	for i, column := range columns {
		if i > 0 {
			cols.WriteString(",")
			params.WriteString(",")
		}
		cols.WriteString(column)
		params.WriteString(fmt.Sprintf("$%d", i+1))
	}
	return cols.String(), params.String()
}

func runInTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := session.Database(ctx).Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			_ = tx.Rollback()
			panic(err)
		}
	}()
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
