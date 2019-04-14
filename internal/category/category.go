package category

import (
	"context"
	"database/sql"
	"godiscourse/internal/durable"
	"time"
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
	Update(ctx context.Context, id string, p *Params) (*Model, error)
	GetByID(ctx context.Context, id string) (*Model, error)
	GetAll(ctx context.Context) ([]*Model, error)
}

type Category struct {
	db *durable.Database
}

func New(db *durable.Database) *Category {
	return &Category{db: db}
}
