package admin

import (
	"satellity/internal/durable"

	"github.com/dimfeld/httptreemux"
)

// RegisterAdminRoutes register admin routes
func RegisterAdminRoutes(database *durable.Database, router *httptreemux.Group) {
	registerAdminUser(database, router)
	registerAdminCategory(database, router)
	registerAdminTopic(database, router)
}
