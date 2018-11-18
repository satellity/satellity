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

type categoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func registerAdminCategory(router *httptreemux.TreeMux) {
	impl := &adminCategoryImpl{}
	router.POST("/admin/categories", impl.create)
	router.GET("/admin/categories", impl.index)
}

func (impl *adminCategoryImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	category, err := models.CreateCategory(r.Context(), body.Name, body.Description)
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
