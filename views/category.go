package views

import (
	"net/http"

	"github.com/godiscourse/godiscourse/models"
)

// CategoryView is the response body of category
type CategoryView struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func buildCategory(category *models.Category) CategoryView {
	return CategoryView{
		Type:        "category",
		Name:        category.Name,
		Description: category.Description,
	}
}

// RenderCategory response a category
func RenderCategory(w http.ResponseWriter, r *http.Request, category *models.Category) {
	RenderResponse(w, r, buildCategory(category))
}

func RenderCategories(w http.ResponseWriter, r *http.Request, categories []*models.Category) {
	views := make([]CategoryView, len(categories))
	for i, c := range categories {
		views[i] = buildCategory(c)
	}
	RenderResponse(w, r, views)
}
