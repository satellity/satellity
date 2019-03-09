package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/godiscourse/godiscourse/api/durable"
	"github.com/godiscourse/godiscourse/api/middleware"
	"github.com/godiscourse/godiscourse/api/models"
	"github.com/godiscourse/godiscourse/api/session"
	"github.com/godiscourse/godiscourse/api/views"
)

type topicImpl struct {
	database *durable.Database
}

type topicRequest struct {
	Title      string `json:"title"`
	Body       string `json:"body"`
	CategoryID string `json:"category_id"`
}

func registerTopic(database *durable.Database, router *httptreemux.TreeMux) {
	impl := &topicImpl{database: database}

	router.POST("/topics", impl.create)
	router.POST("/topics/:id", impl.update)
	router.GET("/topics", impl.index)
	router.GET("/topics/:id", impl.show)
}

func (impl *topicImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body topicRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	ctx := models.WrapContext(r.Context(), impl.database)
	if topic, err := middleware.CurrentUser(r).CreateTopic(ctx, body.Title, body.Body, body.CategoryID); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopic(w, r, topic)
	}
}

func (impl *topicImpl) update(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body topicRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	ctx := models.WrapContext(r.Context(), impl.database)
	if topic, err := middleware.CurrentUser(r).UpdateTopic(ctx, params["id"], body.Title, body.Body, body.CategoryID); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopic(w, r, topic)
	}
}

func (impl *topicImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	ctx := models.WrapContext(r.Context(), impl.database)
	if topic, err := models.ReadTopic(ctx, params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if topic == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else {
		views.RenderTopic(w, r, topic)
	}
}

func (impl *topicImpl) index(w http.ResponseWriter, r *http.Request, params map[string]string) {
	ctx := models.WrapContext(r.Context(), impl.database)
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	if topics, err := models.ReadTopics(ctx, offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}
