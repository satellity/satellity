package views

import (
	"net/http"
	"time"

	"godiscourse/internal/model"
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

func buildTopic(topic *model.Topic) TopicView {
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
	// if topic.User != nil {
	// 	view.User = buildUser(topic.User)
	// }
	// if topic.Category != nil {
	// 	view.Category = buildCategory(topic.Category)
	// }
	return view
}

// RenderTopic response a topic
func RenderTopic(w http.ResponseWriter, r *http.Request, topic *model.Topic) {
	RenderResponse(w, r, buildTopic(topic))
}

// RenderTopics response a bundle of topics
func RenderTopics(w http.ResponseWriter, r *http.Request, topics []*model.Topic) {
	topicViews := make([]TopicView, len(topics))
	for i, t := range topics {
		topicViews[i] = buildTopic(t)
	}
	RenderResponse(w, r, topicViews)
}
