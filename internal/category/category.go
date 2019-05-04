package category

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"godiscourse/internal/durable"
	"godiscourse/internal/session"

	"github.com/gofrs/uuid"
)

type Model struct {
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

type Params struct {
	Name        string
	Alias       string
	Description string
	Position    int64
}

type CategoryDatastore interface {
	Create(ctx context.Context, p *Params) (*Model, error)
	Find(ctx context.Context, tx *sql.Tx, id string) (*Model, error)
	Update(ctx context.Context, id string, p *Params) (*Model, error)
	GetByID(ctx context.Context, id string) (*Model, error)
	GetAll(ctx context.Context) ([]*Model, error)
	GetCategorySet(ctx context.Context, tx *sql.Tx) (map[string]*Model, error)
}

type Category struct {
	db *durable.Database
}

func New(db *durable.Database) *Category {
	return &Category{db: db}
}

func (c *Category) Create(ctx context.Context, p *Params) (*Model, error) {
	alias, name := strings.TrimSpace(p.Alias), strings.TrimSpace(p.Name)
	description := strings.TrimSpace(p.Description)
	if len(name) < 1 {
		return nil, session.BadDataError(ctx)
	}
	if alias == "" {
		alias = p.Name
	}

	t := time.Now()
	category := &Model{
		CategoryID:  uuid.Must(uuid.NewV4()).String(),
		Name:        name,
		Alias:       alias,
		Description: description,
		TopicsCount: 0,
		LastTopicID: sql.NullString{String: "", Valid: false},
		Position:    p.Position,
		CreatedAt:   t,
		UpdatedAt:   t,
	}

	cols, params := durable.PrepareColumnsWithValues(categoryColumns)
	err := c.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		if p.Position == 0 {
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

func (c *Category) Find(ctx context.Context, tx *sql.Tx, id string) (*Model, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}

	row := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM categories WHERE category_id=$1", strings.Join(categoryColumns, ",")), id)
	cat, err := categoryFromRows(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return cat, err
}

func (c *Category) Update(ctx context.Context, id string, p *Params) (*Model, error) {
	alias, name := strings.TrimSpace(p.Alias), strings.TrimSpace(p.Name)
	description := strings.TrimSpace(p.Description)
	if len(alias) < 1 && len(name) < 1 {
		return nil, session.BadDataError(ctx)
	}

	var category *Model
	err := c.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		category, err := c.Find(ctx, tx, id)
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
		category.Position = p.Position
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

func (c *Category) GetByID(ctx context.Context, id string) (*Model, error) {
	return &Model{}, nil
}

func (c *Category) GetAll(ctx context.Context) ([]*Model, error) {
	var categories []*Model
	err := c.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		categories, err = readCategories(ctx, tx)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return categories, nil
}

func (c *Category) GetCategorySet(ctx context.Context, tx *sql.Tx) (map[string]*Model, error) {
	categories, err := readCategories(ctx, tx)
	if err != nil {
		return nil, err
	}
	set := make(map[string]*Model, 0)
	for _, c := range categories {
		set[c.CategoryID] = c
	}
	return set, nil
}
