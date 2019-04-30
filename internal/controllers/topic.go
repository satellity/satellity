package controllers

import (
	"encoding/json"
	"godiscourse/internal/durable"
	"godiscourse/internal/middleware"
	"godiscourse/internal/models"
	"godiscourse/internal/session"
	"godiscourse/internal/views"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
)

type topicImpl struct {
	database *durable.Database
}

type topicRequest struct {
	Title      string `json:"title"`
	Body       string `json:"body"`
	CategoryID string `json:"category_id"`
}

func registerTopic(database *durable.Database, router *httptreemux.TreeMux) {
	impl := &topicImpl{database: database}

	router.POST("/topics", impl.create)
	router.POST("/topics/:id", impl.update)
	router.GET("/topics", impl.index)
	router.GET("/topics/:id", impl.show)
	router.POST("/topics/:id/like", impl.like)
	router.POST("/topics/:id/unlike", impl.unlike)
	router.POST("/topics/:id/bookmark", impl.bookmark)
	router.POST("/topics/:id/abandon", impl.abandon)
}

func (impl *topicImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body topicRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	mctx := models.WrapContext(r.Context(), impl.database)
	if topic, err := middleware.CurrentUser(r).CreateTopic(mctx, body.Title, body.Body, body.CategoryID); err != nil {
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
	mctx := models.WrapContext(r.Context(), impl.database)
	if topic, err := middleware.CurrentUser(r).UpdateTopic(mctx, params["id"], body.Title, body.Body, body.CategoryID); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopic(w, r, topic)
	}
}

func (impl *topicImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	if topic, err := models.ReadTopicWithUser(mctx, params["id"], middleware.CurrentUser(r)); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if topic == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else {
		views.RenderTopic(w, r, topic)
	}
}

func (impl *topicImpl) index(w http.ResponseWriter, r *http.Request, params map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	if topics, err := models.ReadTopics(mctx, offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}

func (impl *topicImpl) like(w http.ResponseWriter, r *http.Request, params map[string]string) {
	impl.action(w, r, params["id"], models.TopicUserActionLiked, true)
}

func (impl *topicImpl) unlike(w http.ResponseWriter, r *http.Request, params map[string]string) {
	impl.action(w, r, params["id"], models.TopicUserActionLiked, false)
}

func (impl *topicImpl) bookmark(w http.ResponseWriter, r *http.Request, params map[string]string) {
	impl.action(w, r, params["id"], models.TopicUserActionBookmarked, true)
}

func (impl *topicImpl) abandon(w http.ResponseWriter, r *http.Request, params map[string]string) {
	impl.action(w, r, params["id"], models.TopicUserActionBookmarked, false)
}

func (impl *topicImpl) action(w http.ResponseWriter, r *http.Request, id, action string, state bool) {
	mctx := models.WrapContext(r.Context(), impl.database)
	topic, err := models.ReadTopic(mctx, id)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	} else if topic == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
		return
	}
	if topic, err := topic.ActiondBy(mctx, middleware.CurrentUser(r), action, state); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopic(w, r, topic)
	}
}
