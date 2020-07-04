package admin

import (
	"net/http"
	"satellity/internal/middlewares"
	"satellity/internal/models"
	"satellity/internal/views"
	"time"

	"github.com/dimfeld/httptreemux"
)

type topicImpl struct{}

func registerAdminTopic(router *httptreemux.Group) {
	impl := &topicImpl{}

	router.DELETE("/topics/:id", impl.destroy)
	router.GET("/topics", impl.index)
}

func (impl *topicImpl) destroy(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if err := middlewares.CurrentUser(r).DeleteTopic(r.Context(), params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderBlankResponse(w, r)
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
