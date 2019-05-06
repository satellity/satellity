package views

import (
	"net/http"
	"time"

	"godiscourse/internal/models"
)

// TopicView is the response body of topic
type TopicView struct {
	Type          string       `json:"type"`
	TopicID       string       `json:"topic_id"`
	ShortID       string       `json:"short_id"`
	Title         string       `json:"title"`
	Body          string       `json:"body"`
	UserID        string       `json:"user_id"`
	CategoryID    string       `json:"category_id"`
	Score         int          `json:"score"`
	CommentsCount int64        `json:"comments_count"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	User          UserView     `json:"user"`
	Category      CategoryView `json:"category"`
}

func buildTopic(topic *models.Topic) TopicView {
	view := TopicView{
		Type:          "topic",
		TopicID:       topic.TopicID,
		ShortID:       topic.ShortID,
		Title:         topic.Title,
		Body:          topic.Body,
		UserID:        topic.UserID,
		CategoryID:    topic.CategoryID,
		Score:         topic.Score,
		CommentsCount: topic.CommentsCount,
		CreatedAt:     topic.CreatedAt,
		UpdatedAt:     topic.UpdatedAt,
	}
	view.User = buildUser(&topic.User)
	view.Category = buildCategory(&topic.Category)
	return view
}

// RenderTopic response a topic
func RenderTopic(w http.ResponseWriter, r *http.Request, topic *models.Topic) {
	RenderResponse(w, r, buildTopic(topic))
}

// RenderTopics response a bundle of topics
func RenderTopics(w http.ResponseWriter, r *http.Request, topics []*models.Topic) {
	topicViews := make([]TopicView, len(topics))
	for i, t := range topics {
		topicViews[i] = buildTopic(t)
	}
	RenderResponse(w, r, topicViews)
}
