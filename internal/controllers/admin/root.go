package admin

import (
	"github.com/dimfeld/httptreemux"
)

// RegisterAdminRoutes register admin routes
func RegisterAdminRoutes(router *httptreemux.Group) {
	api := router.NewGroup("/admin")
	registerAdminUser(api)
	registerAdminCategory(api)
	registerAdminTopic(api)
	registerAdminComment(api)
	registerAdminProduct(api)
}
