package views

import (
	"net/http"
	"time"

	"github.com/godiscourse/godiscourse/models"
)

type UserView struct {
	Type      string    `json:"type"`
	UserId    string    `json:"user_id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AccountView struct {
	UserView
	SessionId string `json:"session_id"`
}

func buildUser(user *models.User) UserView {
	return UserView{
		Type:      "user",
		UserId:    user.UserId,
		Email:     user.Email.String,
		Username:  user.Username,
		Nickname:  user.Nickname,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func RenderUser(w http.ResponseWriter, r *http.Request, user *models.User) {
	RenderResponse(w, r, buildUser(user))
}

func RenderUsers(w http.ResponseWriter, r *http.Request, users []*models.User) {
	userViews := make([]UserView, len(users))
	for i, user := range users {
		userViews[i] = buildUser(user)
	}
	RenderResponse(w, r, userViews)
}

func RenderAccount(w http.ResponseWriter, r *http.Request, user *models.User) {
	accountView := AccountView{
		UserView:  buildUser(user),
		SessionId: user.SessionId,
	}
	RenderResponse(w, r, accountView)
}
