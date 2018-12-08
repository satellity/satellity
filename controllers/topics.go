package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/godiscourse/godiscourse/middleware"
	"github.com/godiscourse/godiscourse/models"
	"github.com/godiscourse/godiscourse/session"
	"github.com/godiscourse/godiscourse/views"
)

type topicRequest struct {
	Title      string `json:"title"`
	Body       string `json:"body"`
	CategoryID string `json:"category_id"`
}

type topicImpl struct{}

func registerTopic(router *httptreemux.TreeMux) {
	impl := &topicImpl{}

	router.POST("/topics", impl.create)
	router.POST("/topics/:id", impl.update)
	router.GET("/topics", impl.index)
	router.GET("/topics/:id", impl.show)
	router.GET("/user/topics", impl.topics)
}

func (impl *topicImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body topicRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	if topic, err := middleware.CurrentUser(r).CreateTopic(r.Context(), body.Title, body.Body, body.CategoryID); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopic(w, r, topic)
	}
}

func (impl *topicImpl) update(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body topicRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	if topic, err := middleware.CurrentUser(r).UpdateTopic(r.Context(), params["id"], body.Title, body.Body); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopic(w, r, topic)
	}
}

func (impl *topicImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if topic, err := models.ReadTopic(r.Context(), params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if topic == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else {
		views.RenderTopic(w, r, topic)
	}
}

func (impl *topicImpl) index(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	if topics, err := models.ReadTopics(r.Context(), offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}

func (impl *topicImpl) topics(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	if topics, err := middleware.CurrentUser(r).ReadTopics(r.Context(), offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}
