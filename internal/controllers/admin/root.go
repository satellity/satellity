package admin

import (
	"godiscourse/internal/durable"

	"github.com/dimfeld/httptreemux"
)

// RegisterAdminRoutes register admin routes
func RegisterAdminRoutes(database *durable.Database, router *httptreemux.TreeMux) {
	registerAdminCategory(database, router)
	registerAdminTopic(database, router)
}
