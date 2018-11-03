package views

import (
	"net/http"

	"github.com/godiscourse/godiscourse/models"
)

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

func RenderCategory(w http.ResponseWriter, r *http.Request, category *models.Category) {
	RenderResponse(w, r, buildCategory(category))
}
