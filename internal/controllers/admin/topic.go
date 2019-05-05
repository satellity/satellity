package admin

import (
	"godiscourse/internal/engine"
	"godiscourse/internal/views"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
)

type topicImpl struct {
	engine engine.Engine
}

func RegisterAdminTopic(e engine.Engine, router *httptreemux.TreeMux) {
	impl := &topicImpl{engine: e}

	router.GET("/admin/topics", impl.index)
}

func (impl *topicImpl) index(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	if topics, err := impl.engine.GetTopicsByOffset(r.Context(), offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}
