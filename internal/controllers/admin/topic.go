package admin

import (
	"net/http"
	"satellity/internal/durable"
	"satellity/internal/middlewares"
	"satellity/internal/models"
	"satellity/internal/views"
	"time"

	"github.com/dimfeld/httptreemux"
)

type topicImpl struct {
	database *durable.Database
}

func registerAdminTopic(database *durable.Database, router *httptreemux.Group) {
	impl := &topicImpl{database: database}

	router.DELETE("/topics/:id", impl.destroy)
	router.GET("/topics", impl.index)
}

func (impl *topicImpl) destroy(w http.ResponseWriter, r *http.Request, params map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	if err := middlewares.CurrentUser(r).DeleteTopic(mctx, params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderBlankResponse(w, r)
	}
}

func (impl *topicImpl) index(w http.ResponseWriter, r *http.Request, params map[string]string) {
	ctx := models.WrapContext(r.Context(), impl.database)
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	if topics, err := models.ReadTopics(ctx, offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderTopics(w, r, topics)
	}
}
