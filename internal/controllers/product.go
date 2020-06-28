package controllers

import (
	"net/http"
	"satellity/internal/models"
	"satellity/internal/session"
	"satellity/internal/views"

	"github.com/dimfeld/httptreemux"
)

type productImpl struct{}

func registerProduct(router *httptreemux.Group) {
	impl := &productImpl{}

	router.GET("/products", impl.index)
	router.GET("/products/:id", impl.show)
}

func (impl *productImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if products, err := models.FindProducts(r.Context()); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderProducts(w, r, products)
	}
}

func (impl *productImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if product, err := models.FindProduct(r.Context(), params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if product == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else {
		views.RenderProduct(w, r, product)
	}
}
