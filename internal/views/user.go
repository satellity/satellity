package views

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"time"

	"godiscourse/internal/user"
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

func buildUser(u *user.Data) UserView {
	return UserView{
		Type:      "user",
		UserID:    u.UserID,
		Nickname:  u.Name(),
		Biography: u.Biography,
		AvatarURL: fmt.Sprintf("https://www.gravatar.com/avatar/%x?s=180&d=wavatar", md5.Sum([]byte(strings.ToLower(u.Email.String)))),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// RenderUser response a user
func RenderUser(w http.ResponseWriter, r *http.Request, u *user.Data) {
	RenderResponse(w, r, buildUser(u))
}

// RenderUsers response a bundle of users
func RenderUsers(w http.ResponseWriter, r *http.Request, users []*user.Data) {
	userViews := make([]UserView, len(users))
	for i, user := range users {
		userViews[i] = buildUser(user)
	}
	RenderResponse(w, r, userViews)
}

// RenderAccount response
func RenderAccount(w http.ResponseWriter, r *http.Request, u *user.Data) {
	accountView := AccountView{
		UserView:  buildUser(u),
		Username:  u.Username,
		Email:     u.Email.String,
		SessionID: u.SessionID,
		Role:      u.Role(),
	}
	RenderResponse(w, r, accountView)
}
