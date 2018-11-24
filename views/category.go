package views

import (
	"net/http"
	"time"

	"github.com/godiscourse/godiscourse/models"
)

// CategoryView is the response body of category
type CategoryView struct {
	Type        string    `json:"type"`
	CategoryID  string    `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TopicsCount int       `json:"topics_count"`
	LastTopicID string    `json:"last_topic_id"`
	Position    int       `json:"position"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func buildCategory(category *models.Category) CategoryView {
	return CategoryView{
		Type:        "category",
		CategoryID:  category.CategoryID,
		Name:        category.Name,
		Description: category.Description,
		TopicsCount: category.TopicsCount,
		LastTopicID: category.LastTopicID.String,
		Position:    category.Position,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

// RenderCategory response a category
func RenderCategory(w http.ResponseWriter, r *http.Request, category *models.Category) {
	RenderResponse(w, r, buildCategory(category))
}

// RenderCategories response sevaral categories
func RenderCategories(w http.ResponseWriter, r *http.Request, categories []*models.Category) {
	views := make([]CategoryView, len(categories))
	for i, c := range categories {
		views[i] = buildCategory(c)
	}
	RenderResponse(w, r, views)
}
