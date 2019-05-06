package controllers

import (
	"encoding/json"
	"godiscourse/internal/engine"
	"godiscourse/internal/middleware"
	"godiscourse/internal/models"
	"godiscourse/internal/session"
	"godiscourse/internal/views"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
)

type topicImpl struct {
	poster engine.Poster
}

type topicRequest struct {
	Title      string `json:"title"`
	Body       string `json:"body"`
	CategoryID string `json:"category_id"`
}

func registerTopic(p engine.Poster, router *httptreemux.TreeMux) {
	impl := &topicImpl{poster: p}

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
	if t, err := impl.poster.CreateTopic(r.Context(), u.UserID, &models.TopicInfo{
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

	if t, err := impl.poster.UpdateTopic(r.Context(), params["id"], &models.TopicInfo{
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
	if t, err := impl.poster.GetTopicByID(r.Context(), params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if t == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else {
		views.RenderTopic(w, r, t)
	}
}

func (impl *topicImpl) index(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	if topics, err := impl.poster.GetTopicsByOffset(r.Context(), offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}
