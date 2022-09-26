package views

import (
	"net/http"
	"satellity/internal/models"
	"time"
)

type GistView struct {
	Type      string           `json:"type"`
	GistID    string           `json:"gist_id"`
	Title     string           `json:"title"`
	Author    string           `json:"author"`
	Link      string           `json:"link"`
	PublishAt time.Time        `json:"publish_at"`
	Source    SourceSimpleView `json:"source"`
}

func buildGist(g *models.Gist) GistView {
	view := GistView{
		Type:      "gist",
		GistID:    g.GistID,
		Title:     g.Title,
		Author:    g.Author,
		Link:      g.Link,
		PublishAt: g.PublishAt,
	}
	if g.Source != nil {
		view.Source = buildSimpleSource(g.Source)
	}
	return view
}

func RenderGists(w http.ResponseWriter, r *http.Request, gists []*models.Gist) {
	views := make([]GistView, len(gists))
	for i, gist := range gists {
		views[i] = buildGist(gist)
	}
	RenderResponse(w, r, views)
}
