package admin

import "github.com/dimfeld/httptreemux"

func RegisterAdminRoutes(router *httptreemux.TreeMux) {
	registerAdminCategory(router)
}
