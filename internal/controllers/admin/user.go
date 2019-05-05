package admin

import (
	"godiscourse/internal/engine"
	"godiscourse/internal/views"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
)

type userImpl struct {
	admin engine.Admin
}

func RegisterAdminUser(a engine.Admin, router *httptreemux.TreeMux) {
	impl := &userImpl{admin: a}

	router.GET("/admin/users", impl.index)
}

func (impl *userImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	users, err := impl.admin.GetUsersByOffset(r.Context(), offset)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderUsers(w, r, users)
}
