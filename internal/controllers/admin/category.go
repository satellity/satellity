package admin

import (
	"encoding/json"
	"godiscourse/internal/durable"
	"godiscourse/internal/models"
	"godiscourse/internal/session"
	"godiscourse/internal/views"
	"net/http"

	"github.com/dimfeld/httptreemux"
)

type adminCategoryImpl struct {
	database *durable.Database
}

type categoryRequest struct {
	Name        string `json:"name"`
	Alias       string `json:"alias"`
	Description string `json:"description"`
	Position    int64  `json:"position"`
}

func registerAdminCategory(database *durable.Database, router *httptreemux.TreeMux) {
	impl := &adminCategoryImpl{database: database}

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
	ctx := models.WrapContext(r.Context(), impl.database)
	category, err := models.CreateCategory(ctx, body.Name, body.Alias, body.Description, body.Position)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderCategory(w, r, category)
}

func (impl *adminCategoryImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	ctx := models.WrapContext(r.Context(), impl.database)
	categories, err := models.ReadAllCategories(ctx)
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
	ctx := models.WrapContext(r.Context(), impl.database)
	category, err := models.UpdateCategory(ctx, params["id"], body.Name, body.Alias, body.Description, body.Position)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderCategory(w, r, category)
}

func (impl *adminCategoryImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	ctx := models.WrapContext(r.Context(), impl.database)
	category, err := models.ReadCategory(ctx, params["id"])
	if err != nil {
		views.RenderErrorResponse(w, r, err)
		return
	}
	views.RenderCategory(w, r, category)
}
