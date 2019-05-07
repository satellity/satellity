package controllers

import (
	"encoding/json"
	"godiscourse/internal/engine"
	"godiscourse/internal/middleware"
	"godiscourse/internal/models"
	"godiscourse/internal/session"
	"godiscourse/internal/user"
	"godiscourse/internal/views"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
)

type userImpl struct {
	user   user.UserDatastore
	poster engine.Poster
}

type userRequest struct {
	Code          string `json:"code"`
	SessionSecret string `json:"session_secret"`
	Nickname      string `json:"nickname"`
	Biography     string `json:"biography"`
}

func RegisterUser(u user.UserDatastore, p engine.Poster, router *httptreemux.TreeMux) {
	impl := &userImpl{
		user:   u,
		poster: p,
	}

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
	if user, err := impl.user.CreateGithubUser(r.Context(), body.Code, body.SessionSecret); err != nil {
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

	result := middleware.CurrentUser(r)
	if err := impl.poster.UpdateUser(r.Context(), result, &models.UserInfo{
		Nickname:  body.Nickname,
		Biography: body.Biography,
	}); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderAccount(w, r, result)
	}
}

func (impl *userImpl) current(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	views.RenderAccount(w, r, middleware.CurrentUser(r))
}

func (impl *userImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if user, err := impl.poster.GetUserByID(r.Context(), params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if user == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else {
		views.RenderUser(w, r, user)
	}
}

func (impl *userImpl) topics(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	user, err := impl.poster.GetUserByID(r.Context(), params["id"])

	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if user == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else if topics, err := impl.poster.GetTopicsByUserID(r.Context(), user.UserID, offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}
