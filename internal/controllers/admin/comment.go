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

type commentImpl struct {
	database *durable.Database
}

func registerAdminComment(database *durable.Database, router *httptreemux.Group) {
	impl := &commentImpl{database: database}

	router.GET("/comments", impl.index)
	router.DELETE("/comments/:id", impl.destroy)
}

func (impl *commentImpl) destroy(w http.ResponseWriter, r *http.Request, params map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	if err := middlewares.CurrentUser(r).DeleteComment(mctx, params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderBlankResponse(w, r)
	}
}

func (impl *commentImpl) index(w http.ResponseWriter, r *http.Request, params map[string]string) {
	ctx := models.WrapContext(r.Context(), impl.database)
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	if comments, err := models.ReadComments(ctx, offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderComments(w, r, comments)
	}
}
