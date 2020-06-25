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

type topicImpl struct{}

type topicRequest struct {
	Title      string `json:"title"`
	Body       string `json:"body"`
	TopicType  string `json:"topic_type"`
	CategoryID string `json:"category_id"`
	Draft      bool   `json:"draft"`
}

func registerTopic(router *httptreemux.Group) {
	impl := &topicImpl{}

	router.POST("/topics", impl.create)
	router.POST("/topics/:id", impl.update)
	router.POST("/topics/:id/like", impl.like)
	router.POST("/topics/:id/unlike", impl.unlike)
	router.POST("/topics/:id/bookmark", impl.bookmark)
	router.POST("/topics/:id/unsave", impl.unsave)
	router.GET("/topics", impl.index)
	router.GET("/topics/draft", impl.draft)
	router.GET("/topics/:id", impl.show)
}

func (impl *topicImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body topicRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	if topic, err := middlewares.CurrentUser(r).CreateTopic(r.Context(), body.Title, body.Body, body.TopicType, body.CategoryID, body.Draft); err != nil {
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
	if topic, err := middlewares.CurrentUser(r).UpdateTopic(r.Context(), params["id"], body.Title, body.Body, body.TopicType, body.CategoryID, body.Draft); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopic(w, r, topic)
	}
}

func (impl *topicImpl) draft(w http.ResponseWriter, r *http.Request, params map[string]string) {
	user := middlewares.CurrentUser(r)
	if user == nil {
		views.RenderErrorResponse(w, r, session.AuthorizationError(r.Context()))
		return
	}
	if topic, err := user.DraftTopic(r.Context()); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if topic == nil {
		views.RenderBlankResponse(w, r)
	} else {
		views.RenderTopic(w, r, topic)
	}
}

func (impl *topicImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if topic, err := models.ReadTopicWithRelation(r.Context(), params["id"], middlewares.CurrentUser(r)); err != nil {
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

func (impl *topicImpl) like(w http.ResponseWriter, r *http.Request, params map[string]string) {
	impl.action(w, r, params["id"], models.TopicUserActionLiked, true)
}

func (impl *topicImpl) unlike(w http.ResponseWriter, r *http.Request, params map[string]string) {
	impl.action(w, r, params["id"], models.TopicUserActionLiked, false)
}

func (impl *topicImpl) bookmark(w http.ResponseWriter, r *http.Request, params map[string]string) {
	impl.action(w, r, params["id"], models.TopicUserActionBookmarked, true)
}

func (impl *topicImpl) unsave(w http.ResponseWriter, r *http.Request, params map[string]string) {
	impl.action(w, r, params["id"], models.TopicUserActionBookmarked, false)
}

func (impl *topicImpl) action(w http.ResponseWriter, r *http.Request, id, action string, state bool) {
	if topic, err := models.ReadTopic(r.Context(), id); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if topic == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else if topic, err = topic.ActiondBy(r.Context(), middlewares.CurrentUser(r), action, state); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopic(w, r, topic)
	}
}
