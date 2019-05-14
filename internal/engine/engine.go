package engine

import (
	"context"
	"godiscourse/internal/models"
	"time"
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
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, userID string, info *models.UserInfo) error

	GetCategoryByID(ctx context.Context, id string) (*models.Category, error)
	GetAllCategories(ctx context.Context) ([]*models.Category, error)

	CreateTopic(ctx context.Context, userID string, t *models.TopicInfo) (*models.Topic, error)
	UpdateTopic(ctx context.Context, id string, t *models.TopicInfo) (*models.Topic, error)
	GetTopicByID(ctx context.Context, id string) (*models.Topic, error)
	GetTopicsByUserID(ctx context.Context, userID string, offset time.Time) ([]*models.Topic, error)
	GetTopicsByCategoryID(ctx context.Context, categoryID string, offset time.Time) ([]*models.Topic, error)
	GetTopicsByOffset(ctx context.Context, offset time.Time) ([]*models.Topic, error)

	CreateComment(ctx context.Context, c *models.CommentInfo) (*models.Comment, error)
	UpdateComment(ctx context.Context, c *models.CommentInfo) (*models.Comment, error)
	DeleteComment(ctx context.Context, id, userID string) error
	GetCommentsByTopicID(ctx context.Context, topicID string, offset time.Time) ([]*models.Comment, error)
}

type Admin interface {
	GetUsersByOffset(ctx context.Context, offset time.Time) ([]*models.User, error)

	CreateCategory(ctx context.Context, c *models.CategoryRequest) (*models.Category, error)
	UpdateCategory(ctx context.Context, id string, c *models.CategoryRequest) (*models.Category, error)
}
