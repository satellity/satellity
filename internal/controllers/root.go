package controllers

import (
	"godiscourse/internal/configs"
	"godiscourse/internal/controllers/admin"
	"godiscourse/internal/engine"
	"godiscourse/internal/views"
	"net/http"
	"runtime"

	"github.com/dimfeld/httptreemux"
)

func Register(eng engine.Engine, router *httptreemux.TreeMux) {
	router.GET("/_hc", health)
	registerHanders(router)

	registerCategory(eng, router)
	registerComment(eng, router)
	registerTopic(eng, router)

	admin.RegisterAdminUser(eng, router)
	admin.RegisterAdminCategory(eng, router)
	admin.RegisterAdminTopic(eng, router)
}

func health(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	views.RenderResponse(w, r, map[string]string{
		"build":      configs.BuildVersion + "-" + runtime.Version(),
		"developers": configs.HTTPResourceHost,
	})
}
