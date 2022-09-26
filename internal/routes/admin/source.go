package admin

import (
	"net/http"
	"satellity/internal/models"
	"satellity/internal/views"

	"github.com/dimfeld/httptreemux"
)

type sourceImpl struct{}

func registerAdminSource(router *httptreemux.Group) {
	impl := &sourceImpl{}

	router.GET("/sources", impl.index)
}

func (impl *sourceImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	sources, err := models.ReadSources(r.Context())
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderSources(w, r, sources)
	}
}
