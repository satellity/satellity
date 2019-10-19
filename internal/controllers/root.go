package controllers

import (
	"net/http"
	"runtime"
	"satellity/internal/configs"
	"satellity/internal/controllers/admin"
	"satellity/internal/durable"
	"satellity/internal/session"
	"satellity/internal/views"

	"github.com/dimfeld/httptreemux"
)

// RegisterRoutes register all routes
func RegisterRoutes(database *durable.Database, router *httptreemux.TreeMux) {
	api := router.NewGroup("/api")
	api.GET("/_hc", health)
	registerUser(database, api)
	registerCategory(database, api)
	registerTopic(database, api)
	registerComment(database, api)
	registerVerification(database, api)
	admin.RegisterAdminRoutes(database, api)
}

func health(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	views.RenderResponse(w, r, map[string]string{
		"build": configs.BuildVersion + "-" + runtime.Version(),
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
