package views

import (
	"net/http"
	"time"

	"github.com/godiscourse/godiscourse/models"
)

// TopicView is the response body of topic
type TopicView struct {
	Type      string    `json:"type"`
	TopicID   string    `json:"topic_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func buildTopic(topic *models.Topic) TopicView {
	return TopicView{
		Type:      "topic",
		TopicID:   topic.TopicID,
		Title:     topic.Title,
		Body:      topic.Body,
		Score:     topic.Score,
		CreatedAt: topic.CreatedAt,
		UpdatedAt: topic.UpdatedAt,
	}
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
