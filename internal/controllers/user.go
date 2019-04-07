package controllers

import (
	"encoding/json"
	"godiscourse/internal/middleware"
	"godiscourse/internal/session"
	"godiscourse/internal/user"
	"godiscourse/internal/views"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
)

type userImpl struct {
	repo *user.UserDatastore
}

type userRequest struct {
	Code          string `json:"code"`
	SessionSecret string `json:"session_secret"`
	Nickname      string `json:"nickname"`
	Biography     string `json:"biography"`
}

func RegisterUser(r *user.User, router *httptreemux.TreeMux) {
	impl := &userImpl{repo: r}

	router.POST("/oauth/:provider", impl.oauth)
	router.POST("/me", impl.update)
	router.GET("/me", impl.current)
	router.GET("/users/:id", impl.show)
	router.GET("/users/:id/topics", impl.topics)
}

func (impl *userImpl) oauth(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body userRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	if user, err := impl.repo.CreateGithubUser(r.Context(), body.Code, body.SessionSecret); err != nil {
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

	user := middleware.CurrentUser(r)
	if err := user.Update(r.Context(), body.Nickname, body.Biography); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderAccount(w, r, user)
	}
}

func (impl *userImpl) current(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	views.RenderAccount(w, r, middleware.CurrentUser(r))
}

func (impl *userImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if user, err := impl.repo.GetByID(r.Context(), params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if user == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else {
		views.RenderUser(w, r, user)
	}
}

func (impl *userImpl) topics(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	user, err := impl.repo.GetByID(ctx, params["id"])

	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if user == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else if topics, err := user.ReadTopics(ctx, offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}
