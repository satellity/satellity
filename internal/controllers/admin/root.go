package admin

import (
	"satellity/internal/durable"

	"github.com/dimfeld/httptreemux"
)

// RegisterAdminRoutes register admin routes
func RegisterAdminRoutes(database *durable.Database, router *httptreemux.Group) {
	api := router.NewGroup("/admin")
	registerAdminUser(database, api)
	registerAdminCategory(database, api)
	registerAdminTopic(database, api)
}
