package controllers

import (
	"net/http"

	"github.com/dimfeld/httptreemux"
)

type userImpl struct{}

func registerUser(router *httptreemux.TreeMux) {
	impl := &userImpl{}
	router.POST("/users", impl.create)
}

func (impl *userImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
}
