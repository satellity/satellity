package controllers

import (
	"net/http"

	"github.com/dimfeld/httptreemux"
)

type categoryImpl struct{}

func registerCategory(router *httptreemux.TreeMux) {
	impl := &categoryImpl{}
	router.GET("/categories", impl.index)
}

func (impl *categoryImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
}
