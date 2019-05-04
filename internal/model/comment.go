package model

import (
	"time"
)

type Comment struct {
	CommentID string    `sql:"comment_id,pk"`
	Body      string    `sql:"body"`
	TopicID   string    `sql:"topic_id"`
	UserID    string    `sql:"user_id"`
	Score     int       `sql:"score,notnull"`
	CreatedAt time.Time `sql:"created_at"`
	UpdatedAt time.Time `sql:"updated_at"`
}

type CommentInfo struct {
	CommentID string
	TopicID   string
	UserID    string
	Body      string
}
