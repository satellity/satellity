package engine

import (
	"context"
	"time"

	"godiscourse/internal/model"
)

// Topic related CONST
const (
	minTitleSize = 3
	LIMIT        = 50
)

// Comment related CONST
const (
	minCommentBodySize = 6
)

type Engine interface {
	Poster
	Admin
}

type Poster interface {
	GetCategoryByID(ctx context.Context, id string) (*model.Category, error)
	GetAllCategories(ctx context.Context) ([]*model.Category, error)

	CreateTopic(ctx context.Context, userID string, t *model.TopicInfo) (*model.Topic, error)
	UpdateTopic(ctx context.Context, id string, t *model.TopicInfo) (*model.Topic, error)
	GetTopicByID(ctx context.Context, id string) (*model.Topic, error)
	GetTopicByUserID(ctx context.Context, userID string, offset time.Time) ([]*model.Topic, error)
	GetTopicsByCategoryID(ctx context.Context, categoryID string, offset time.Time) ([]*model.Topic, error)
	GetTopicsByOffset(ctx context.Context, offset time.Time) ([]*model.Topic, error)

	CreateComment(ctx context.Context, c *model.CommentInfo) (*model.Comment, error)
	UpdateComment(ctx context.Context, c *model.CommentInfo) (*model.Comment, error)
	DeleteComment(ctx context.Context, id, userID string) error
	GetCommentsByTopicID(ctx context.Context, topicID string, offset time.Time) ([]*model.Comment, error)
}

type Admin interface {
	CreateCategory(ctx context.Context, c *model.CategoryInfo) (*model.Category, error)
	UpdateCategory(ctx context.Context, id string, c *model.CategoryInfo) (*model.Category, error)
}
