package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/godiscourse/godiscourse/api/middleware"
	"github.com/godiscourse/godiscourse/api/models"
	"github.com/godiscourse/godiscourse/api/session"
	"github.com/godiscourse/godiscourse/api/views"
)

type userRequest struct {
	Code          string `json:"code"`
	SessionSecret string `json:"session_secret"`
	Nickname      string `json:"nickname"`
	Biography     string `json:"biography"`
}

type userImpl struct{}

func registerUser(router *httptreemux.TreeMux) {
	impl := &userImpl{}

	router.POST("/oauth/:provider", impl.oauth)
	router.POST("/me", impl.update)
	router.GET("/me", impl.self)
	router.GET("/users/:id", impl.show)
	router.GET("/users/:id/topics", impl.topics)
}

func (impl *userImpl) oauth(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body userRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	user, err := models.CreateGithubUser(r.Context(), body.Code, body.SessionSecret)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderAccount(w, r, user)
}

func (impl *userImpl) update(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body userRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	current := middleware.CurrentUser(r)
	err := current.UpdateProfile(r.Context(), body.Nickname, body.Biography)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderAccount(w, r, current)
}

func (impl *userImpl) self(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	views.RenderAccount(w, r, middleware.CurrentUser(r))
}

func (impl *userImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if user, err := models.ReadUser(r.Context(), params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if user == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else {
		views.RenderUser(w, r, user)
	}
}

func (impl *userImpl) topics(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	user, err := models.ReadUser(r.Context(), params["id"])
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if user == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else if topics, err := middleware.CurrentUser(r).ReadTopics(r.Context(), offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}
