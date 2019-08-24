package controllers

import (
	"encoding/json"
	"net/http"
	"satellity/internal/durable"
	"satellity/internal/middlewares"
	"satellity/internal/models"
	"satellity/internal/session"
	"satellity/internal/views"
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
	router.DELETE("/comments/:id", impl.destory)
	router.GET("/topics/:id/comments", impl.comments)
}

func (impl *commentImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body commentRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	mctx := models.WrapContext(r.Context(), impl.database)
	if comment, err := middlewares.CurrentUser(r).CreateComment(mctx, body.TopicID, body.Body); err != nil {
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

	mctx := models.WrapContext(r.Context(), impl.database)
	if comment, err := middlewares.CurrentUser(r).CreateComment(mctx, params["id"], body.Body); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderComment(w, r, comment)
	}
}

func (impl *commentImpl) destory(w http.ResponseWriter, r *http.Request, params map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	if err := middlewares.CurrentUser(r).DeleteComment(mctx, params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderBlankResponse(w, r)
	}
}

func (impl *commentImpl) comments(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	mctx := models.WrapContext(r.Context(), impl.database)
	if topic, err := models.ReadTopic(mctx, params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if topic == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else if comments, err := topic.ReadComments(mctx, offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderComments(w, r, comments)
	}
}
