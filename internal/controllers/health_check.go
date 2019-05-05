package controllers

import (
	"godiscourse/internal/configs"
	"godiscourse/internal/views"
	"net/http"
	"runtime"

	"github.com/dimfeld/httptreemux"
)

func healthCheck(router *httptreemux.TreeMux) {
	router.GET("/_hc", health)
}

func health(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	views.RenderResponse(w, r, map[string]string{
		"build":      configs.BuildVersion + "-" + runtime.Version(),
		"developers": "https://live.godiscourse.com",
	})
}
