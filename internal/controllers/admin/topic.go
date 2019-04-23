package admin

import (
	"godiscourse/internal/topic"
	"godiscourse/internal/views"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
)

type topicImpl struct {
	topic topic.TopicDatastore
}

func RegisterAdminTopic(t topic.TopicDatastore, router *httptreemux.TreeMux) {
	impl := &topicImpl{topic: t}

	router.GET("/admin/topics", impl.index)
}

func (impl *topicImpl) index(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	if topics, err := impl.topic.GetByOffset(r.Context(), offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}
