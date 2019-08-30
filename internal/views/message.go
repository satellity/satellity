package views

import (
	"net/http"
	"satellity/internal/models"
	"time"
)

// MessageView is the response body of message
type MessageView struct {
	Type      string    `json:"type"`
	MessageID string    `json:"message_id"`
	Body      string    `json:"body"`
	GroupID   string    `json:"group_id"`
	UserID    string    `json:"user_id"`
	ParentID  string    `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      UserView  `json:"user"`
}

func buildMessage(message *models.Message) MessageView {
	return MessageView{
		Type:      "message",
		MessageID: message.MessageID,
		Body:      message.Body,
		GroupID:   message.GroupID,
		UserID:    message.UserID,
		ParentID:  message.ParentID,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
		User:      buildUser(message.User),
	}
}

// RenderMessage response a message view
func RenderMessage(w http.ResponseWriter, r *http.Request, message *models.Message) {
	RenderResponse(w, r, buildMessage(message))
}

// RenderMessages response a bundle of messages view
func RenderMessages(w http.ResponseWriter, r *http.Request, messages []*models.Message) {
	views := make([]MessageView, len(messages))
	for i, message := range messages {
		views[i] = buildMessage(message)
	}
	RenderResponse(w, r, views)
}
