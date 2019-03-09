package admin

import (
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/godiscourse/godiscourse/api/durable"
	"github.com/godiscourse/godiscourse/api/models"
	"github.com/godiscourse/godiscourse/api/views"
)

type userImpl struct {
	database *durable.Database
}

func registerAdminUser(database *durable.Database, router *httptreemux.TreeMux) {
	impl := &userImpl{database: database}

	router.GET("/admin/users", impl.index)
}

func (impl *userImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	ctx := models.WrapContext(r.Context(), impl.database)
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	users, err := models.ReadUsers(ctx, offset)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderUsers(w, r, users)
}
