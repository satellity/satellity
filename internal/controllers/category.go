package controllers

import (
	"godiscourse/internal/engine"
	"godiscourse/internal/session"
	"godiscourse/internal/views"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
)

type categoryImpl struct {
	poster engine.Poster
}

func registerCategory(p engine.Poster, router *httptreemux.TreeMux) {
	impl := &categoryImpl{poster: p}

	router.GET("/categories", impl.index)
	router.GET("/categories/:id/topics", impl.topics)
}

func (impl *categoryImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	categories, err := impl.poster.GetAllCategories(r.Context())
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}

	views.RenderCategories(w, r, categories)
}

func (impl *categoryImpl) topics(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))

	category, err := impl.poster.GetCategoryByID(r.Context(), params["id"])
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if category == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else if topics, err := impl.poster.GetTopicsByCategoryID(r.Context(), category.CategoryID, offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}
