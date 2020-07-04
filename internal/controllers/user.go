package controllers

import (
	"encoding/json"
	"net/http"
	"satellity/internal/middlewares"
	"satellity/internal/models"
	"satellity/internal/session"
	"satellity/internal/views"
	"time"

	"github.com/dimfeld/httptreemux"
)

type userImpl struct{}

type userRequest struct {
	Code          string `json:"code"`
	SessionSecret string `json:"session_secret"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Nickname      string `json:"nickname"`
	Avatar        string `json:"avatar"`
	Biography     string `json:"biography"`
}

func registerUser(router *httptreemux.Group) {
	impl := &userImpl{}

	router.POST("/oauth/:provider", impl.oauth)
	router.POST("/sessions", impl.create)
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
	if user, err := models.CreateGithubUser(r.Context(), body.Code, body.SessionSecret); err != nil {
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
	if user, err := models.CreateSession(r.Context(), body.Email, body.Password, body.SessionSecret); err != nil {
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
	current := middlewares.CurrentUser(r)
	if err := current.UpdateProfile(r.Context(), body.Nickname, body.Biography, body.Avatar); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderAccount(w, r, current)
	}
}

func (impl *userImpl) current(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	views.RenderAccount(w, r, middlewares.CurrentUser(r))
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
	} else if topics, err := user.ReadTopics(r.Context(), offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}
