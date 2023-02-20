package routes

import (
	"net/http"
	"satellity/internal/models"
	"satellity/internal/views"

	"github.com/dimfeld/httptreemux"
)

type assetImpl struct{}

func registerAsset(router *httptreemux.Group) {
	impl := &assetImpl{}

	router.GET("/assets", impl.index)
}

func (impl *assetImpl) index(w http.ResponseWriter, r *http.Request, params map[string]string) {
	assets, err := models.ReadAssetsWithRatios(r.Context())
	if err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderAssets(w, r, assets)
	}
}
