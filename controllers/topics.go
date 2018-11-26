package controllers

import (
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/godiscourse/godiscourse/middleware"
	"github.com/godiscourse/godiscourse/models"
	"github.com/godiscourse/godiscourse/views"
)

type topicImpl struct{}

func registerTopic(router *httptreemux.TreeMux) {
	impl := &topicImpl{}

	router.GET("/topics/:id", impl.show)
	router.GET("/topics", impl.index)
	router.GET("/user/topics", impl.topics)
}

func (impl *topicImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if topic, err := models.ReadTopic(r.Context(), params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
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
