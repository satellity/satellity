package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func RegisterRoutes(router *httprouter.Router) {
	router.GET("/_hc", health)

	registerUser(router)
}

func health(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	view.RenderBlankResponse(w, r)
}
