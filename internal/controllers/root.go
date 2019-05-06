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

func Register(engine engine.Engine, router *httptreemux.TreeMux) {
	router.GET("/_hc", health)
	registerHanders(router)

	registerCategory(engine, router)
	registerComment(engine, router)
	registerTopic(engine, router)

	admin.RegisterAdminUser(engine, router)
	admin.RegisterAdminCategory(engine, router)
	admin.RegisterAdminTopic(engine, router)
}

func health(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	views.RenderResponse(w, r, map[string]string{
		"build":      configs.BuildVersion + "-" + runtime.Version(),
		"developers": configs.HTTPResourceHost,
	})
}
