package main

import (
	"context"
	"database/sql"
	"fmt"
	"godiscourse/internal/configs"
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
	"go.uber.org/zap"
)

func startHTTP(db *sql.DB, logger *zap.Logger) error {
	database := durable.WrapDatabase(db)
	router := httptreemux.New()
	controllers.RegisterHanders(router)
	controllers.RegisterRoutes(database, router)

	handler := middleware.Authenticate(database, router)
	handler = middleware.Constraint(handler)
	handler = middleware.Context(handler, render.New())
	handler = middleware.State(handler)
	handler = middleware.Logger(handler, durable.NewLogger(logger))
	handler = handlers.ProxyHeaders(handler)

	return gracehttp.Serve(&http.Server{Addr: fmt.Sprintf(":%d", configs.HTTPListenPort), Handler: handler})
}

func main() {
	db := durable.OpenDatabaseClient(context.Background())
	defer db.Close()

	logger, err := zap.NewDevelopment()
	if configs.Environment == "production" {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatal(err)
	}

	if err := startHTTP(db, logger); err != nil {
		log.Panicln(err)
	}
}
