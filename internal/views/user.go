package views

import (
	"net/http"
	"satellity/internal/models"
	"time"
)

// UserView is the response body of user
type UserView struct {
	Type      string    `json:"type"`
	UserID    string    `json:"user_id"`
	Nickname  string    `json:"nickname"`
	Biography string    `json:"biography"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AccountView is the response body of a sign in user
type AccountView struct {
	UserView
	Username  string `json:"username"`
	Email     string `json:"email"`
	SessionID string `json:"session_id"`
	Role      string `json:"role"`
}

func buildUser(user *models.User) UserView {
	return UserView{
		Type:      "user",
		UserID:    user.UserID,
		Nickname:  user.Name(),
		Biography: user.Biography,
		AvatarURL: user.GetAvatar(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// RenderUser response a user
func RenderUser(w http.ResponseWriter, r *http.Request, user *models.User) {
	RenderResponse(w, r, buildUser(user))
}

// RenderUsers response a bundle of users
func RenderUsers(w http.ResponseWriter, r *http.Request, users []*models.User) {
	userViews := make([]UserView, len(users))
	for i, user := range users {
		userViews[i] = buildUser(user)
	}
	RenderResponse(w, r, userViews)
}

// RenderAccount response
func RenderAccount(w http.ResponseWriter, r *http.Request, user *models.User) {
	accountView := AccountView{
		UserView:  buildUser(user),
		Username:  user.Username,
		Email:     user.Email.String,
		SessionID: user.SessionID,
		Role:      user.GetRole(),
	}
	RenderResponse(w, r, accountView)
}
