package controllers

import (
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/godiscourse/godiscourse/models"
	"github.com/godiscourse/godiscourse/views"
)

type categoryImpl struct{}

func registerCategory(router *httptreemux.TreeMux) {
	impl := &categoryImpl{}

	router.GET("/categories", impl.index)
}

func (impl *categoryImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	categories, err := models.ReadCategories(r.Context())
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}

	views.RenderCategories(w, r, categories)
}
