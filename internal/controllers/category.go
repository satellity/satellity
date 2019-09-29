package controllers

import (
	"net/http"
	"satellity/internal/durable"
	"satellity/internal/models"
	"satellity/internal/session"
	"satellity/internal/views"
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
	mctx := models.WrapContext(r.Context(), impl.database)
	categories, err := models.ReadAllCategories(mctx)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderCategories(w, r, categories)
	}
}

func (impl *categoryImpl) topics(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	mctx := models.WrapContext(r.Context(), impl.database)
	category, err := models.ReadCategory(mctx, params["id"])
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if category == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else if topics, err := category.ReadTopics(mctx, offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}
