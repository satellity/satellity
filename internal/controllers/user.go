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

type userImpl struct {
	database *durable.Database
}

type userRequest struct {
	Code          string `json:"code"`
	SessionSecret string `json:"session_secret"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Nickname      string `json:"nickname"`
	Biography     string `json:"biography"`
}

func registerUser(database *durable.Database, router *httptreemux.TreeMux) {
	impl := &userImpl{database: database}

	router.POST("/oauth/:provider", impl.oauth)
	router.POST("/sessions", impl.create)
	router.POST("/me", impl.update)
	router.GET("/me", impl.current)
	router.GET("/users/:id", impl.show)
	router.GET("/users/:id/topics", impl.topics)
	router.GET("/users/:id/groups", impl.groups)
	router.GET("/user/groups", impl.relatedGroups)
}

func (impl *userImpl) oauth(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body userRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	mctx := models.WrapContext(r.Context(), impl.database)
	if user, err := models.CreateGithubUser(mctx, body.Code, body.SessionSecret); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderAccount(w, r, user)
	}
}

func (impl *userImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body userRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	mctx := models.WrapContext(r.Context(), impl.database)
	if user, err := models.CreateSession(mctx, body.Email, body.Password, body.SessionSecret); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderAccount(w, r, user)
	}
}

func (impl *userImpl) update(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body userRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	mctx := models.WrapContext(r.Context(), impl.database)
	current := middlewares.CurrentUser(r)
	if err := current.UpdateProfile(mctx, body.Nickname, body.Biography); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderAccount(w, r, current)
	}
}

func (impl *userImpl) current(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	views.RenderAccount(w, r, middlewares.CurrentUser(r))
}

func (impl *userImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	if user, err := models.ReadUser(mctx, params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if user == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else {
		views.RenderUser(w, r, user)
	}
}

func (impl *userImpl) topics(w http.ResponseWriter, r *http.Request, params map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	user, err := models.ReadUser(mctx, params["id"])
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if user == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else if topics, err := user.ReadTopics(mctx, offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}

func (impl *userImpl) groups(w http.ResponseWriter, r *http.Request, params map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	groups, err := models.ReadGroupsByUser(mctx, params["id"])
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderGroups(w, r, groups)
	}
}

func (impl *userImpl) relatedGroups(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)

	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	groups, err := middlewares.CurrentUser(r).RelatedGroups(mctx, limit)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderGroups(w, r, groups)
	}
}
