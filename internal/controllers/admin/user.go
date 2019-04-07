package admin

import (
	"godiscourse/internal/user"
	"godiscourse/internal/views"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
)

type userImpl struct {
	repo *user.UserDatastore
}

func RegisterAdminUser(r *user.User, router *httptreemux.TreeMux) {
	impl := &userImpl{repo: r}

	router.GET("/admin/users", impl.index)
}

func (impl *userImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	users, err := impl.repo.GetByOffset(r.Context(), offset)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderUsers(w, r, users)
}
