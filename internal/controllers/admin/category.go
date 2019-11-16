package admin

import (
	"encoding/json"
	"net/http"
	"satellity/internal/durable"
	"satellity/internal/models"
	"satellity/internal/session"
	"satellity/internal/views"

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

func registerAdminCategory(database *durable.Database, router *httptreemux.Group) {
	impl := &adminCategoryImpl{database: database}

	router.POST("/categories", impl.create)
	router.POST("/categories/:id", impl.update)
	router.GET("/categories", impl.index)
	router.GET("/categories/:id", impl.show)
}

func (impl *adminCategoryImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	mctx := models.WrapContext(r.Context(), impl.database)
	category, err := models.CreateCategory(mctx, body.Name, body.Alias, body.Description, body.Position)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderCategory(w, r, category)
	}
}

func (impl *adminCategoryImpl) index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	categories, err := models.ReadAllCategories(mctx)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderCategories(w, r, categories)
	}
}

func (impl *adminCategoryImpl) update(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body categoryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	mctx := models.WrapContext(r.Context(), impl.database)
	category, err := models.UpdateCategory(mctx, params["id"], body.Name, body.Alias, body.Description, body.Position)
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderCategory(w, r, category)
	}
}

func (impl *adminCategoryImpl) show(w http.ResponseWriter, r *http.Request, params map[string]string) {
	mctx := models.WrapContext(r.Context(), impl.database)
	category, err := models.ReadCategory(mctx, params["id"])
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderCategory(w, r, category)
	}
}
