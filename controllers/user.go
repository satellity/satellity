package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type userImpl struct{}

func registerUser(router *httprouter.Router) {
	impl := &userImpl{}
	router.POST("/users", impl.create)
}

func (impl *userImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
}
