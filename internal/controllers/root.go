package controllers

import (
	"godiscourse/internal/controllers/admin"
	"godiscourse/internal/engine"

	"github.com/dimfeld/httptreemux"
)

func Register(engine engine.Engine, router *httptreemux.TreeMux) {
	healthCheck(router)
	registerHanders(router)

	registerCategory(engine, router)
	registerComment(engine, router)
	registerTopic(engine, router)

	admin.RegisterAdminCategory(engine, router)
	admin.RegisterAdminTopic(engine, router)
}
