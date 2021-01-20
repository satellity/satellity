package controllers

import (
	"net/http"
	"satellity/internal/models"
	"satellity/internal/session"
	"satellity/internal/views"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/dimfeld/httptreemux"
	"github.com/gofrs/uuid"
)

type productImpl struct{}

func registerProduct(router *httptreemux.Group) {
	impl := &productImpl{}

	router.GET("/products", impl.index)
	router.GET("/products/:id/relationships", impl.relationships)
	router.GET("/products/:id", impl.show)
}

func (impl *productImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	query := r.URL.Query().Get("q")
	var products []*models.Product
	var err error
	if query != "" {
		products, err = models.SearchProducts(r.Context(), query)
	} else {
		products, err = models.FindProducts(r.Context())
	}
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderProducts(w, r, products)
	}
}

func (impl *productImpl) relationships(w http.ResponseWriter, r *http.Request, params map[string]string) {
	if products, err := models.RelatedProducts(r.Context(), params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderProducts(w, r, products)
	}
}

func (impl *productImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	subs := strings.Split(params["id"], "-")
	id, err := uuid.FromBytes(base58.Decode(subs[len(subs)-1]))
	if err != nil {
		views.RenderErrorResponse(w, r, session.ServerError(r.Context(), err))
		return
	}
	if product, err := models.FindProduct(r.Context(), id.String()); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if product == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else {
		views.RenderProduct(w, r, product)
	}
}
