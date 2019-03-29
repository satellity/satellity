package controllers

import (
	"godiscourse/internal/durable"
	"godiscourse/internal/models"
	"godiscourse/internal/session"
	"godiscourse/internal/views"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
)

type categoryImpl struct {
	database *durable.Database
}

func registerCategory(database *durable.Database, router *httptreemux.TreeMux) {
	impl := &categoryImpl{database: database}

	router.GET("/categories", impl.index)
	router.GET("/categories/:id/topics", impl.topics)
}

func (impl *categoryImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	ctx := models.WrapContext(r.Context(), impl.database)
	categories, err := models.ReadAllCategories(ctx)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}

	views.RenderCategories(w, r, categories)
}

func (impl *categoryImpl) topics(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	ctx := models.WrapContext(r.Context(), impl.database)
	category, err := models.ReadCategory(ctx, params["id"])
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if category == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else if topics, err := category.ReadTopics(ctx, offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}
