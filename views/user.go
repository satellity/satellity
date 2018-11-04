package views

import (
	"net/http"
	"time"

	"github.com/godiscourse/godiscourse/models"
)

// UserView is the response body of user
type UserView struct {
	Type      string    `json:"type"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AccountView is the response body of a sign in user
type AccountView struct {
	UserView
	SessionID string `json:"session_id"`
	Role      string `json:"role"`
}

func buildUser(user *models.User) UserView {
	return UserView{
		Type:      "user",
		UserID:    user.UserID,
		Email:     user.Email.String,
		Username:  user.Username,
		Nickname:  user.Nickname,
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
		SessionID: user.SessionID,
		Role:      user.Role(),
	}
	RenderResponse(w, r, accountView)
}
