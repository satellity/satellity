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

type messageRequest struct {
	Body string `json:"body"`
}

type messageImpl struct {
	database *durable.Database
}

func registerMessage(database *durable.Database, router *httptreemux.TreeMux) {
	impl := &messageImpl{database: database}
	router.POST("/groups/:id/messages", impl.create)
	router.POST("/messages/:id", impl.update)
	router.POST("/messages/:id/delete", impl.destroy)
	router.GET("/groups/:id/messages", impl.index)
}

func (impl *messageImpl) create(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body messageRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}

	mctx := models.WrapContext(r.Context(), impl.database)
	message, err := middleware.CurrentUser(r).CreateMessage(mctx, params["id"], body.Body)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderMessage(w, r, message)
	}
}

func (impl *messageImpl) update(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body messageRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}

	mctx := models.WrapContext(r.Context(), impl.database)
	message, err := middleware.CurrentUser(r).UpdateMessage(mctx, params["id"], body.Body)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderMessage(w, r, message)
	}
}

func (impl *messageImpl) destroy(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body messageRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}

	mctx := models.WrapContext(r.Context(), impl.database)
	err := middleware.CurrentUser(r).DeleteMessage(mctx, params["id"])
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderBlankResponse(w, r)
	}
}

func (impl *messageImpl) index(w http.ResponseWriter, r *http.Request, params map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))

	group, err := models.ReadGroup(mctx, params["id"], middleware.CurrentUser(r))
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if group == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else if messages, err := group.ReadMessages(mctx, offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderMessages(w, r, messages)
	}
}
