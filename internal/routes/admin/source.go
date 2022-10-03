package admin

import (
	"net/http"
	"satellity/internal/models"
	"satellity/internal/session"
	"satellity/internal/views"

	"github.com/dimfeld/httptreemux"
)

type sourceImpl struct{}

func registerAdminSource(router *httptreemux.Group) {
	impl := &sourceImpl{}

	router.GET("/sources", impl.index)
	router.DELETE("/sources/:id", impl.destroy)
}

func (impl *sourceImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	sources, err := models.ReadSources(r.Context())
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderSources(w, r, sources)
	}
}

func (impl *sourceImpl) destroy(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if source, err := models.ReadSource(r.Context(), params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if source == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else if err = source.Delete(r.Context()); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderBlankResponse(w, r)
	}
}
