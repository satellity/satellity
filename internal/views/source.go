package views

import "satellity/internal/models"

type SourceView struct {
	Author  string `json:"author"`
	LogoURL string `json:"logo_url"`
}

func buildSource(source *models.Source) *SourceView {
	return &SourceView{
		Author:  source.Author,
		LogoURL: source.LogoURL,
	}
}
