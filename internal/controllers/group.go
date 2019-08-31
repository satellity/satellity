package controllers

import (
	"encoding/json"
	"net/http"
	"satellity/internal/durable"
	"satellity/internal/middlewares"
	"satellity/internal/models"
	"satellity/internal/session"
	"satellity/internal/views"
	"strconv"
	"time"

	"github.com/dimfeld/httptreemux"
)

type groupImpl struct {
	database *durable.Database
}

type groupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Cover       string `json:"cover"`
}

func registerGroup(database *durable.Database, router *httptreemux.TreeMux) {
	impl := &groupImpl{database: database}

	router.POST("/groups", impl.create)
	router.POST("/groups/:id", impl.update)
	router.POST("/groups/:id/join", impl.join)
	router.POST("/groups/:id/exit", impl.exit)
	router.GET("/groups/:id/participants", impl.participants)
	router.GET("/groups/:id/messages", impl.messages)

	router.GET("/groups", impl.index)
	router.GET("/groups/:id", impl.show)
}

func (impl *groupImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body groupRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	mctx := models.WrapContext(r.Context(), impl.database)
	if group, err := middlewares.CurrentUser(r).CreateGroup(mctx, body.Name, body.Description, body.Cover); err != nil {
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
	if group, err := middlewares.CurrentUser(r).UpdateGroup(mctx, params["id"], body.Name, body.Description, body.Cover); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if group == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else {
		views.RenderGroup(w, r, group)
	}
}

func (impl *groupImpl) join(w http.ResponseWriter, r *http.Request, params map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	if group, err := middlewares.CurrentUser(r).JoinGroup(mctx, params["id"], models.ParticipantRoleMember); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if group == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else {
		views.RenderGroup(w, r, group)
	}
}

func (impl *groupImpl) exit(w http.ResponseWriter, r *http.Request, params map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	if group, err := middlewares.CurrentUser(r).ExitGroup(mctx, params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if group == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else {
		views.RenderGroup(w, r, group)
	}
}

func (impl *groupImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if groups, err := models.ReadGroups(mctx, offset, limit); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderGroups(w, r, groups)
	}
}

func (impl *groupImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	if group, err := models.ReadGroup(mctx, params["id"], middlewares.CurrentUser(r)); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if group == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else {
		views.RenderGroup(w, r, group)
	}
}

func (impl *groupImpl) participants(w http.ResponseWriter, r *http.Request, params map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))

	current := middlewares.CurrentUser(r)
	if group, err := models.ReadGroup(mctx, params["id"], current); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if group == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else if users, err := group.Participants(mctx, current, offset, r.URL.Query().Get("limit")); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderUsers(w, r, users)
	}
}

func (impl *groupImpl) messages(w http.ResponseWriter, r *http.Request, params map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))

	group, err := models.ReadGroup(mctx, params["id"], middlewares.CurrentUser(r))
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
