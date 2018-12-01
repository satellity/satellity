package views

import (
	"net/http"
	"time"

	"github.com/godiscourse/godiscourse/models"
)

// TopicView is the response body of topic
type TopicView struct {
	Type          string       `json:"type"`
	TopicID       string       `json:"topic_id"`
	Title         string       `json:"title"`
	Body          string       `json:"body"`
	Score         int          `json:"score"`
	CommentsCount int          `json:"comments_count"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	User          UserView     `json:"user"`
	Category      CategoryView `json:"category"`
}

func buildTopic(topic *models.Topic) TopicView {
	view := TopicView{
		Type:          "topic",
		TopicID:       topic.TopicID,
		Title:         topic.Title,
		Body:          topic.Body,
		Score:         topic.Score,
		CommentsCount: topic.CommentsCount,
		CreatedAt:     topic.CreatedAt,
		UpdatedAt:     topic.UpdatedAt,
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
