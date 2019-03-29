package main

import (
	"context"
	"database/sql"
	"fmt"
	"godiscourse/internal/config"
	"godiscourse/internal/controllers"
	"godiscourse/internal/durable"
	"godiscourse/internal/middleware"
	"log"
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/handlers"
	_ "github.com/lib/pq"
	"github.com/unrolled/render"
)

func startHTTP(db *sql.DB) error {
	database := durable.WrapDatabase(db)
	router := httptreemux.New()
	controllers.RegisterHanders(router)
	controllers.RegisterRoutes(database, router)

	handler := middleware.Authenticate(database, router)
	handler = middleware.Constraint(handler)
	handler = middleware.Context(handler, render.New())
	handler = middleware.State(handler)
	handler = middleware.Logger(handler, durable.NewLogger())
	handler = handlers.ProxyHeaders(handler)

	return gracehttp.Serve(&http.Server{Addr: fmt.Sprintf(":%d", config.HTTPListenPort), Handler: handler})
}

func main() {
	db := durable.OpenDatabaseClient(context.Background())
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Panicln(err)
	}

	if err := startHTTP(db); err != nil {
		log.Panicln(err)
	}
}
