package controllers

import (
	"encoding/json"
	"godiscourse/internal/middleware"
	"godiscourse/internal/session"
	"godiscourse/internal/topic"
	"godiscourse/internal/user"
	"godiscourse/internal/views"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
)

type topicImpl struct {
	repo      topic.TopicDatastore
	userStore user.UserDatastore
}

type topicRequest struct {
	Title      string `json:"title"`
	Body       string `json:"body"`
	CategoryID string `json:"category_id"`
}

func RegisterTopic(t topic.TopicDatastore, us user.UserDatastore, router *httptreemux.TreeMux) {
	impl := &topicImpl{
		repo:      t,
		userStore: us,
	}

	router.POST("/topics", impl.create)
	router.POST("/topics/:id", impl.update)
	router.GET("/topics", impl.index)
	router.GET("/topics/:id", impl.show)
}

func (impl *topicImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body topicRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}

	u := middleware.CurrentUser(r)
	if t, err := impl.repo.Create(r.Context(), u.UserID, &topic.Params{
		Title:      body.Title,
		Body:       body.Body,
		CategoryID: body.CategoryID,
	}); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopic(w, r, t)
	}
}

func (impl *topicImpl) update(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body topicRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}

	u := middleware.CurrentUser(r)
	if t, err := impl.repo.Update(r.Context(), u.UserID, params["id"], &topic.Params{
		Title:      body.Title,
		Body:       body.Body,
		CategoryID: body.CategoryID,
	}); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopic(w, r, t)
	}
}

func (impl *topicImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if t, err := impl.repo.GetByID(r.Context(), params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if t == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else {
		views.RenderTopic(w, r, t)
	}
}

func (impl *topicImpl) index(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	if topics, err := impl.repo.GetByOffset(r.Context(), offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}
