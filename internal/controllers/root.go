package controllers

import (
	"godiscourse/internal/configs"
	"godiscourse/internal/controllers/admin"
	"godiscourse/internal/durable"
	"godiscourse/internal/views"
	"net/http"
	"runtime"

	"github.com/dimfeld/httptreemux"
)

// RegisterRoutes register all routes
func RegisterRoutes(database *durable.Database, router *httptreemux.TreeMux) {
	router.GET("/_hc", health)
	registerCategory(database, router)
	registerTopic(database, router)
	registerComment(database, router)
	admin.RegisterAdminRoutes(database, router)
}

func health(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	views.RenderResponse(w, r, map[string]string{
		"build":      configs.BuildVersion + "-" + runtime.Version(),
		"developers": "https://live.godiscourse.com",
	})
}
