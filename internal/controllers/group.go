package controllers

import (
	"encoding/json"
	"godiscourse/internal/durable"
	"godiscourse/internal/middleware"
	"godiscourse/internal/models"
	"godiscourse/internal/session"
	"godiscourse/internal/views"
	"net/http"

	"github.com/dimfeld/httptreemux"
)

type groupImpl struct {
	database *durable.Database
}

type groupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func registerGroup(database *durable.Database, router *httptreemux.TreeMux) {
	impl := &groupImpl{database: database}

	router.POST("/groups", impl.create)
	router.POST("/groups/:id", impl.update)
	router.POST("/groups/:id/join", impl.join)
	router.POST("/groups/:id/exit", impl.exit)
	router.GET("/groups/:id", impl.show)
	router.GET("/groups/:id/participants", impl.participants)
}

func (impl *groupImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body groupRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	mctx := models.WrapContext(r.Context(), impl.database)
	if group, err := middleware.CurrentUser(r).CreateGroup(mctx, body.Name, body.Description); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderGroup(w, r, group)
	}
}

func (impl *groupImpl) update(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body groupRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	mctx := models.WrapContext(r.Context(), impl.database)
	if group, err := middleware.CurrentUser(r).UpdateGroup(mctx, params["id"], body.Name, body.Description); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderGroup(w, r, group)
	}
}

func (impl *groupImpl) join(w http.ResponseWriter, r *http.Request, params map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	if err := middleware.CurrentUser(r).JoinGroup(mctx, params["id"], models.ParticipantRoleMember); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderBlankResponse(w, r)
	}
}

func (impl *groupImpl) exit(w http.ResponseWriter, r *http.Request, params map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	if err := middleware.CurrentUser(r).ExitGroup(mctx, params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderBlankResponse(w, r)
	}
}

func (impl *groupImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	if group, err := models.ReadGroup(mctx, params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderGroup(w, r, group)
	}
}

func (impl *groupImpl) participants(w http.ResponseWriter, r *http.Request, params map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	if group, err := models.ReadGroup(mctx, params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if group == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else if users, err := group.Participants(mctx); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderUsers(w, r, users)
	}
}
