package admin

import (
	"encoding/json"
	"net/http"
	"satellity/internal/middlewares"
	"satellity/internal/session"
	"satellity/internal/views"

	"github.com/dimfeld/httptreemux"
)

type productRequest struct {
	Name   string   `json:"name"`
	Body   string   `json:"body"`
	Cover  string   `json:"cover"`
	Source string   `json:"source"`
	Tags   []string `json:"tags"`
}

type productImpl struct{}

func registerAdminProduct(router *httptreemux.Group) {
	impl := &productImpl{}

	router.POST("/products", impl.create)
	router.POST("/products/:id", impl.update)
}

func (impl *productImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body productRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	if product, err := middlewares.CurrentUser(r).CreateProduct(r.Context(), body.Name, body.Body, body.Cover, body.Source, body.Tags); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderProduct(w, r, product)
	}
}

func (impl *productImpl) update(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body productRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	if product, err := middlewares.CurrentUser(r).UpdateProduct(r.Context(), params["id"], body.Name, body.Body, body.Cover, body.Source, body.Tags); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderProduct(w, r, product)
	}
}
