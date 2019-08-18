package views

import (
	"net/http"
	"satellity/internal/models"
	"time"
)

type GroupView struct {
	Type        string    `json:"type"`
	GroupID     string    `json:"group_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UsersCount  int64     `json:"users_count"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UserView    UserView  `json:"user"`
}

func buildGroup(group *models.Group) GroupView {
	view := GroupView{
		Type:        "group",
		GroupID:     group.GroupID,
		Name:        group.Name,
		Description: group.Description,
		UsersCount:  group.UsersCount,
		Role:        group.GetRole(),
		CreatedAt:   group.CreatedAt,
	}
	if group.User != nil {
		view.UserView = buildUser(group.User)
	}
	return view
}

// RenderGroup response a group view
func RenderGroup(w http.ResponseWriter, r *http.Request, group *models.Group) {
	RenderResponse(w, r, buildGroup(group))
}

// RenderGroups response a bundle of group views
func RenderGroups(w http.ResponseWriter, r *http.Request, groups []*models.Group) {
	views := make([]GroupView, len(groups))
	for i, group := range groups {
		views[i] = buildGroup(group)
	}
	RenderResponse(w, r, views)
}
