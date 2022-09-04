package admin

import (
	"net/http"
	"satellity/internal/models"
	"satellity/internal/views"
	"time"

	"github.com/dimfeld/httptreemux"
)

type userImpl struct{}

func registerAdminUser(router *httptreemux.Group) {
	impl := &userImpl{}

	router.GET("/users", impl.index)
}

func (impl *userImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	users, err := models.ReadUsers(r.Context(), offset)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderUsers(w, r, users)
	}
}
