package views

import (
	"net/http"
	"satellity/internal/models"
	"time"
)

type SourceSimpleView struct {
	Author  string `json:"author"`
	LogoURL string `json:"logo_url"`
	Host    string `json:"host"`
}

type SourceView struct {
	SourceID    string    `json:"source_id"`
	Author      string    `json:"author"`
	Host        string    `json:"host"`
	Link        string    `json:"link"`
	LogoURL     string    `json:"logo_url"`
	Locality    string    `json:"locality"`
	Wreck       int64     `json:"wreck"`
	PublishedAt time.Time `json:"publish_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func buildSimpleSource(s *models.Source) SourceSimpleView {
	return SourceSimpleView{
		Author:  s.Author,
		LogoURL: s.LogoURL,
		Host:    s.Host,
	}
}

func buildSource(s *models.Source) SourceView {
	return SourceView{
		SourceID:    s.SourceID,
		Author:      s.Author,
		Host:        s.Host,
		Link:        s.Link,
		LogoURL:     s.LogoURL,
		Locality:    s.Locality,
		Wreck:       s.Wreck,
		PublishedAt: s.PublishedAt,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func RenderSources(w http.ResponseWriter, r *http.Request, sources []*models.Source) {
	views := make([]SourceView, len(sources))
	for i, s := range sources {
		views[i] = buildSource(s)
	}
	RenderResponse(w, r, views)
}
