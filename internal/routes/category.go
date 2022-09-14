package routes

import (
	"net/http"
	"satellity/internal/models"
	"satellity/internal/session"
	"satellity/internal/views"
	"time"

	"github.com/dimfeld/httptreemux"
)

type categoryImpl struct{}

func registerCategory(router *httptreemux.Group) {
	impl := &categoryImpl{}

	router.GET("/categories", impl.index)
	router.GET("/categories/:id/topics", impl.topics)
}

func (impl *categoryImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	categories, err := models.ReadAllCategories(r.Context())
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderCategories(w, r, categories)
	}
}

func (impl *categoryImpl) topics(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	var category *models.Category
	var err error
	if params["id"] != "latest" {
		category, err = models.ReadCategoryByIDOrName(r.Context(), params["id"])
		if err != nil {
			views.RenderErrorResponse(w, r, err)
			return
		} else if category == nil {
			views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
			return
		}
	}
	if topics, err := models.ReadTopics(r.Context(), offset, category, nil); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}
