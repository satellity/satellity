package admin

import (
	"net/http"
	"satellity/internal/models"
	"satellity/internal/views"
	"strconv"
	"time"

	"github.com/dimfeld/httptreemux"
)

type gistImpl struct{}

func registerAdminGist(router *httptreemux.Group) {
	impl := &gistImpl{}
	router.GET("/gists", impl.index)
}

func (impl *gistImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if gists, err := models.ReadGists(r.Context(), offset, limit); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderGists(w, r, gists)
	}
}
