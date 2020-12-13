package models

import (
	"context"
	"database/sql"
	"fmt"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/lib/pq"
)

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

func categoryFromRows(row durable.Row) (*Category, error) {
	var c Category
	err := row.Scan(&c.CategoryID, &c.Name, &c.Alias, &c.Description, &c.TopicsCount, &c.LastTopicID, &c.Position, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

// CreateCategory create a new category with none blank name and alias, and optional description.
// alias use for human-readable, position for ordering categories
func CreateCategory(ctx context.Context, name, alias, description string, position int64) (*Category, error) {
	alias, name = strings.TrimSpace(alias), strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" {
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

	err := session.Database(ctx).RunInTransaction(ctx, nil, func(tx *sql.Tx) error {
		if position == 0 {
			count, err := categoryCount(ctx, tx)
			if err != nil {
				return err
			}
			category.Position = count
		}
		stmt, err := tx.PrepareContext(ctx, pq.CopyIn("categories", categoryColumns...))
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.ExecContext(ctx, category.values()...)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return category, nil
}

// UpdateCategory update fields of a category
func UpdateCategory(ctx context.Context, id, name, alias, description string, position int64) (*Category, error) {
	alias, name = strings.TrimSpace(alias), strings.TrimSpace(name)
	description = strings.TrimSpace(description)

	var category *Category
	err := session.Database(ctx).RunInTransaction(ctx, nil, func(tx *sql.Tx) error {
		var err error
		category, err = findCategory(ctx, tx, id)
		if err != nil || category == nil {
			return err
		}
		if name != "" {
			category.Name = name
		}
		if alias != "" {
			category.Alias = alias
		}
		category.Description = description
		category.Position = position
		category.UpdatedAt = time.Now()
		cols, posits := durable.PrepareColumnsWithParams([]string{"name", "alias", "description", "position", "updated_at"})
		values := []interface{}{category.Name, category.Alias, category.Description, category.Position, category.UpdatedAt}
		stmt, err := tx.PrepareContext(ctx, fmt.Sprintf("UPDATE categories SET (%s)=(%s) WHERE category_id='%s'", cols, posits, category.CategoryID))
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.ExecContext(ctx, values...)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return category, nil
}

// ReadCategory read a category by ID
func ReadCategory(ctx context.Context, id string) (*Category, error) {
	var category *Category
	err := session.Database(ctx).RunInTransaction(ctx, nil, func(tx *sql.Tx) error {
		var err error
		category, err = findCategory(ctx, tx, id)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return category, nil
}

// ReadCategoryByIDOrName read a category by id or name
func ReadCategoryByIDOrName(ctx context.Context, identity string) (*Category, error) {
	query := fmt.Sprintf("SELECT %s FROM categories WHERE category_id=$1 OR name=$1", strings.Join(categoryColumns, ","))
	row, err := session.Database(ctx).QueryRowContext(ctx, query, identity)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	c, err := categoryFromRows(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, nil
}

// ReadAllCategories read categories order by position
func ReadAllCategories(ctx context.Context) ([]*Category, error) {
	var categories []*Category
	err := session.Database(ctx).RunInTransaction(ctx, nil, func(tx *sql.Tx) error {
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

// emitToCategory update category's info, e.g.: LastTopicID, TopicsCount
func emitToCategory(db *durable.Database, logger *durable.Logger, id string) (*Category, error) {
	ctx := session.WithDatabase(context.Background(), db)
	ctx = session.WithLogger(ctx, logger)
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}
	var category *Category
	err := session.Database(ctx).RunInTransaction(ctx, nil, func(tx *sql.Tx) error {
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
		cols, posits := durable.PrepareColumnsWithParams([]string{"last_topic_id", "topics_count", "updated_at"})
		values := []interface{}{category.LastTopicID, category.TopicsCount, category.UpdatedAt}
		stmt, err := tx.PrepareContext(ctx, fmt.Sprintf("UPDATE categories SET (%s)=(%s) WHERE category_id='%s'", cols, posits, category.CategoryID))
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.ExecContext(ctx, values...)
		return err
	})
	if err != nil {
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
