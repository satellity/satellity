package controllers

import (
	"encoding/json"
	"net/http"
	"satellity/internal/durable"
	"satellity/internal/middlewares"
	"satellity/internal/models"
	"satellity/internal/views"

	"github.com/dimfeld/httptreemux"
)

type messageRequest struct {
	ParentID string `json:"parent_id"`
	Body     string `json:"body"`
}

type messageImpl struct {
	database *durable.Database
}

func registerMessage(database *durable.Database, router *httptreemux.TreeMux) {
	impl := &messageImpl{database: database}
	router.POST("/groups/:id/messages", impl.create)
	router.POST("/messages/:id", impl.update)
	router.DELETE("/messages/:id", impl.destroy)
}

func (impl *messageImpl) create(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body messageRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}

	mctx := models.WrapContext(r.Context(), impl.database)
	message, err := middlewares.CurrentUser(r).CreateMessage(mctx, params["id"], body.Body, body.ParentID)
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
	message, err := middlewares.CurrentUser(r).UpdateMessage(mctx, params["id"], body.Body)
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
	err := middlewares.CurrentUser(r).DeleteMessage(mctx, params["id"])
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderBlankResponse(w, r)
	}
}
