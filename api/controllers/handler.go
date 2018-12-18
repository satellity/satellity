package controllers

import (
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/godiscourse/godiscourse/api/session"
	"github.com/godiscourse/godiscourse/api/views"
)

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
