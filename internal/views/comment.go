package views

import (
	"net/http"
	"satellity/internal/models"
	"time"
)

// CommentView is the response body of comment, which belongs to a topic
type CommentView struct {
	Type      string    `json:"type"`
	CommentID string    `json:"comment_id"`
	Body      string    `json:"body"`
	TopicID   string    `json:"topic_id"`
	UserID    string    `json:"user_id"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      UserView  `json:"user"`
}

func buildComment(comment *models.Comment) CommentView {
	view := CommentView{
		Type:      "comment",
		CommentID: comment.CommentID,
		Body:      comment.Body,
		TopicID:   comment.TopicID,
		UserID:    comment.UserID,
		Score:     comment.Score,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
	if comment.User != nil {
		view.User = buildUser(comment.User)
	}
	return view
}

// RenderComment response single comment
func RenderComment(w http.ResponseWriter, r *http.Request, comment *models.Comment) {
	RenderResponse(w, r, buildComment(comment))
}

// RenderComments response an array of comments
func RenderComments(w http.ResponseWriter, r *http.Request, comments []*models.Comment) {
	views := make([]CommentView, len(comments))
	for i, comment := range comments {
		views[i] = buildComment(comment)
	}
	RenderResponse(w, r, views)
}
