package admin

import (
	"github.com/dimfeld/httptreemux"
	"github.com/godiscourse/godiscourse/api/durable"
)

// RegisterAdminRoutes register admin routes
func RegisterAdminRoutes(database *durable.Database, router *httptreemux.TreeMux) {
	registerAdminUser(database, router)
	registerAdminCategory(database, router)
	registerAdminTopic(database, router)
}
