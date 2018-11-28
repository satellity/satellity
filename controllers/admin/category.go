package admin

import (
	"encoding/json"
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/godiscourse/godiscourse/models"
	"github.com/godiscourse/godiscourse/session"
	"github.com/godiscourse/godiscourse/views"
)

type adminCategoryImpl struct{}

// TODO should add position
type categoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Position    int    `json:"position"`
}

func registerAdminCategory(router *httptreemux.TreeMux) {
	impl := &adminCategoryImpl{}

	router.POST("/admin/categories", impl.create)
	router.GET("/admin/categories", impl.index)
	router.POST("/admin/categories/:id", impl.update)
	router.GET("/admin/categories/:id", impl.show)
}

func (impl *adminCategoryImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	category, err := models.CreateCategory(r.Context(), body.Name, body.Description, body.Position)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderCategory(w, r, category)
}

func (impl *adminCategoryImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	categories, err := models.ReadCategories(r.Context())
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}

	views.RenderCategories(w, r, categories)
}

func (impl *adminCategoryImpl) update(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	category, err := models.UpdateCategory(r.Context(), params["id"], body.Name, body.Description, body.Position)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderCategory(w, r, category)
}

func (impl *adminCategoryImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	category, err := models.ReadCategory(r.Context(), params["id"])
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderCategory(w, r, category)
}
