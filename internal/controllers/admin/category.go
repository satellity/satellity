package admin

import (
	"encoding/json"
	"net/http"
	"satellity/internal/models"
	"satellity/internal/session"
	"satellity/internal/views"

	"github.com/dimfeld/httptreemux"
)

type categoryImpl struct{}

type categoryRequest struct {
	Name        string `json:"name"`
	Alias       string `json:"alias"`
	Description string `json:"description"`
	Position    int64  `json:"position"`
}

func registerAdminCategory(router *httptreemux.Group) {
	impl := &categoryImpl{}

	router.POST("/categories", impl.create)
	router.POST("/categories/:id", impl.update)
	router.GET("/categories", impl.index)
	router.GET("/categories/:id", impl.show)
}

func (impl *categoryImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	category, err := models.CreateCategory(r.Context(), body.Name, body.Alias, body.Description, body.Position)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderCategory(w, r, category)
	}
}

func (impl *categoryImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	categories, err := models.ReadAllCategories(r.Context())
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderCategories(w, r, categories)
	}
}

func (impl *categoryImpl) update(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	category, err := models.UpdateCategory(r.Context(), params["id"], body.Name, body.Alias, body.Description, body.Position)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if category == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else {
		views.RenderCategory(w, r, category)
	}
}

func (impl *categoryImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	category, err := models.ReadCategory(r.Context(), params["id"])
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if category == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else {
		views.RenderCategory(w, r, category)
	}
}
