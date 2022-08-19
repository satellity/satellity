package admin

import (
	"net/http"
	"satellity/internal/middlewares"
	"satellity/internal/models"
	"satellity/internal/session"
	"satellity/internal/views"
	"time"

	"github.com/dimfeld/httptreemux"
)

type commentImpl struct{}

func registerAdminComment(router *httptreemux.Group) {
	impl := &commentImpl{}

	router.GET("/comments", impl.index)
	router.DELETE("/comments/:id", impl.destroy)
}

func (impl *commentImpl) destroy(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if comment, err := models.ReadComment(r.Context(), params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if comment == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else if err = comment.Delete(r.Context(), middlewares.CurrentUser(r)); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderBlankResponse(w, r)
	}
}

func (impl *commentImpl) index(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	if comments, err := models.ReadComments(r.Context(), offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderComments(w, r, comments)
	}
}
