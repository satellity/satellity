package admin

import (
	"godiscourse/internal/user"
	"godiscourse/internal/views"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
)

type userImpl struct {
	user user.UserDatastore
}

func RegisterAdminUser(u user.UserDatastore, router *httptreemux.TreeMux) {
	impl := &userImpl{user: u}

	router.GET("/admin/users", impl.index)
}

func (impl *userImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	users, err := impl.user.GetByOffset(r.Context(), offset)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderUsers(w, r, users)
}
