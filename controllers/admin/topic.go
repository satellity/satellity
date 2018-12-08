package admin

import (
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/godiscourse/godiscourse/models"
	"github.com/godiscourse/godiscourse/views"
)

type topicImpl struct{}

func registerAdminTopic(router *httptreemux.TreeMux) {
	impl := &topicImpl{}

	router.GET("/admin/topics", impl.index)
}

func (impl *topicImpl) index(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	if topics, err := models.ReadTopics(r.Context(), offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}
