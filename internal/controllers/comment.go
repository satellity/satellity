package controllers

import (
	"encoding/json"
	"godiscourse/internal/durable"
	"godiscourse/internal/middleware"
	"godiscourse/internal/models"
	"godiscourse/internal/session"
	"godiscourse/internal/views"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
)

type commentImpl struct {
	database *durable.Database
}

type commentRequest struct {
	TopicID string `json:"topic_id"`
	Body    string `json:"body"`
}

func registerComment(database *durable.Database, router *httptreemux.TreeMux) {
	impl := &commentImpl{database: database}

	router.POST("/comments", impl.create)
	router.POST("/comments/:id", impl.update)
	router.POST("/comments/:id/delete", impl.destory)
	router.GET("/topics/:id/comments", impl.comments)
}

func (impl *commentImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body commentRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	ctx := models.WrapContext(r.Context(), impl.database)
	if comment, err := middleware.CurrentUser(r).CreateComment(ctx, body.TopicID, body.Body); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderComment(w, r, comment)
	}
}

func (impl *commentImpl) update(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body commentRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	ctx := models.WrapContext(r.Context(), impl.database)
	if comment, err := middleware.CurrentUser(r).CreateComment(ctx, params["id"], body.Body); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderComment(w, r, comment)
	}
}

func (impl *commentImpl) destory(w http.ResponseWriter, r *http.Request, params map[string]string) {
	ctx := models.WrapContext(r.Context(), impl.database)
	if err := middleware.CurrentUser(r).DeleteComment(ctx, params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderBlankResponse(w, r)
	}
}

func (impl *commentImpl) comments(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	ctx := models.WrapContext(r.Context(), impl.database)
	if topic, err := models.ReadTopic(ctx, params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if topic == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else if comments, err := topic.ReadComments(ctx, offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderComments(w, r, comments)
	}
}
