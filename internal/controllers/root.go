package controllers

import (
	"net/http"
	"runtime"
	"satellity/internal/configs"
	"satellity/internal/controllers/admin"
	"satellity/internal/session"
	"satellity/internal/views"

	"github.com/dimfeld/httptreemux"
)

// RegisterRoutes register all routes
func RegisterRoutes(router *httptreemux.TreeMux) {
	api := router.NewGroup("/api")

	api.GET("/_hc", health)
	api.GET("/client", client)
	registerUser(api)
	registerCategory(api)
	registerTopic(api)
	registerComment(api)
	registerProduct(api)
	registerVerification(api)
	admin.RegisterAdminRoutes(api)
}

func health(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	views.RenderResponse(w, r, map[string]string{
		"build": configs.BuildVersion + "-" + runtime.Version(),
	})
}

func client(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	config := configs.AppConfig
	views.RenderResponse(w, r, map[string]string{
		"name":               config.Name,
		"github_client_id":   config.Github.ClientID,
		"recaptcha_site_key": config.Recaptcha.SiteKey,
	})
}

// RegisterHanders handle global responses: MethodNotAllowedHandler, NotFoundHandler, PanicHandler
func RegisterHanders(router *httptreemux.TreeMux) {
	router.MethodNotAllowedHandler = func(w http.ResponseWriter, r *http.Request, _ map[string]httptreemux.HandlerFunc) {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	}
	router.NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	}
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, rcv interface{}) {
		err, _ := rcv.(error)
		views.RenderErrorResponse(w, r, session.ServerError(r.Context(), err))
	}
}
