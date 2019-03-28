package admin

import (
	"godiscourse/durable"

	"github.com/dimfeld/httptreemux"
)

// RegisterAdminRoutes register admin routes
func RegisterAdminRoutes(database *durable.Database, router *httptreemux.TreeMux) {
	registerAdminUser(database, router)
	registerAdminCategory(database, router)
	registerAdminTopic(database, router)
}
