package admin

import (
	"net/http"
	"satellity/internal/durable"
	"satellity/internal/models"
	"satellity/internal/views"
	"time"

	"github.com/dimfeld/httptreemux"
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
	} else {
		views.RenderUsers(w, r, users)
	}
}
