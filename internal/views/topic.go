package views

import (
	"net/http"
	"satellity/internal/models"
	"time"
)

// TopicView is the response body of topic
type TopicView struct {
	Type           string       `json:"type"`
	TopicID        string       `json:"topic_id"`
	ShortID        string       `json:"short_id"`
	Title          string       `json:"title"`
	Body           string       `json:"body"`
	TopicType      string       `json:"topic_type"`
	UserID         string       `json:"user_id"`
	CategoryID     string       `json:"category_id"`
	CommentsCount  int64        `json:"comments_count"`
	LikesCount     int64        `json:"likes_count"`
	ViewsCount     int64        `json:"views_count"`
	BookmarksCount int64        `json:"bookmarks_count"`
	IsLikedBy      bool         `json:"is_liked_by"`
	IsBookmarkedBy bool         `json:"is_bookmarked_by"`
	Draft          bool         `json:"draft"`
	Score          int          `json:"score"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	User           UserView     `json:"user"`
	Category       CategoryView `json:"category"`
}

func buildTopic(topic *models.Topic) TopicView {
	view := TopicView{
		Type:           "topic",
		TopicID:        topic.TopicID,
		ShortID:        topic.ShortID,
		Title:          topic.Title,
		Body:           topic.Body,
		TopicType:      topic.TopicType,
		UserID:         topic.UserID,
		CategoryID:     topic.CategoryID,
		IsLikedBy:      topic.IsLikedBy,
		IsBookmarkedBy: topic.IsBookmarkedBy,
		CommentsCount:  topic.CommentsCount,
		LikesCount:     topic.LikesCount,
		ViewsCount:     topic.ViewsCount,
		BookmarksCount: topic.BookmarksCount,
		Draft:          topic.Draft,
		Score:          topic.Score,
		CreatedAt:      topic.CreatedAt,
		UpdatedAt:      topic.UpdatedAt,
	}
	if topic.User != nil {
		view.User = buildUser(topic.User)
	}
	if topic.Category != nil {
		view.Category = buildCategory(topic.Category)
	}
	return view
}

// RenderTopic response a topic
func RenderTopic(w http.ResponseWriter, r *http.Request, topic *models.Topic) {
	RenderResponse(w, r, buildTopic(topic))
}

// RenderTopics response a bundle of topics
func RenderTopics(w http.ResponseWriter, r *http.Request, topics []*models.Topic) {
	topicViews := make([]TopicView, len(topics))
	for i, topic := range topics {
		topicViews[i] = buildTopic(topic)
	}
	RenderResponse(w, r, topicViews)
}
