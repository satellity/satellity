package admin

import (
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/godiscourse/godiscourse/models"
	"github.com/godiscourse/godiscourse/views"
)

type userImpl struct{}

func registerAdminUser(router *httptreemux.TreeMux) {
	impl := &userImpl{}

	router.GET("/admin/users", impl.index)
}

func (impl *userImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	users, err := models.ReadUsers(r.Context(), offset)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderUsers(w, r, users)
}
