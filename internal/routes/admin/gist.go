package admin

import (
	"net/http"
	"satellity/internal/models"
	"satellity/internal/session"
	"satellity/internal/views"
	"time"

	"github.com/dimfeld/httptreemux"
)

type gistImpl struct{}

func registerAdminGist(router *httptreemux.Group) {
	impl := &gistImpl{}
	router.GET("/gists", impl.index)
	router.DELETE("/gists/:id", impl.destroy)
}

func (impl *gistImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	if gists, err := models.ReadAllGists(r.Context(), offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderGists(w, r, gists)
	}
}

func (impl *gistImpl) destroy(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if gist, err := models.ReadGist(r.Context(), params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if gist == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else if err = gist.Delete(r.Context()); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderBlankResponse(w, r)
	}
}
