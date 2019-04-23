package controllers

import (
	"godiscourse/internal/category"
	"godiscourse/internal/session"
	"godiscourse/internal/topic"
	"godiscourse/internal/views"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
)

type categoryImpl struct {
	category category.CategoryDatastore
	topic    topic.TopicDatastore
}

func RegisterCategory(c category.CategoryDatastore, t topic.TopicDatastore, router *httptreemux.TreeMux) {
	impl := &categoryImpl{
		category: c,
		topic:    t,
	}

	router.GET("/categories", impl.index)
	router.GET("/categories/:id/topics", impl.topics)
}

func (impl *categoryImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	categories, err := impl.category.GetAll(r.Context())
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}

	views.RenderCategories(w, r, categories)
}

func (impl *categoryImpl) topics(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))

	category, err := impl.category.GetByID(r.Context(), params["id"])
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if category == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else if topics, err := impl.topic.GetByCategoryID(r.Context(), category, offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}
