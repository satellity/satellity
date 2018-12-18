package controllers

import (
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/godiscourse/godiscourse/api/controllers/admin"
	"github.com/godiscourse/godiscourse/api/views"
)

// RegisterRoutes register all routes
func RegisterRoutes(router *httptreemux.TreeMux) {
	router.GET("/_hc", health)

	registerUser(router)
	registerCategory(router)
	registerTopic(router)
	registerComment(router)
	admin.RegisterAdminRoutes(router)
}

func health(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	views.RenderBlankResponse(w, r)
}
