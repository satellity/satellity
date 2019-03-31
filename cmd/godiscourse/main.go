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
	"os"

	"github.com/dimfeld/httptreemux"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/handlers"
	flags "github.com/jessevdk/go-flags"
	_ "github.com/lib/pq"
	"github.com/unrolled/render"
	"go.uber.org/zap"
)

func startHTTP(db *sql.DB, logger *zap.Logger, port string) error {
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

	return gracehttp.Serve(&http.Server{Addr: fmt.Sprintf(":%s", port), Handler: handler})
}

func main() {
	opts := configs.DefaultOptions()
	if configs.Environment == "production" {
		p := flags.NewParser(&opts, flags.Default)
		if _, err := p.Parse(); err != nil {
			if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
				os.Exit(0)
			} else {
				os.Exit(1)
			}
		}
	}

	db := durable.OpenDatabaseClient(context.Background(), &durable.ConnectionInfo{
		User:     opts.DbUser,
		Password: opts.DbPassword,
		Host:     opts.DbHost,
		Port:     opts.DbPort,
		Name:     opts.DbName,
	})
	defer db.Close()

	logger, err := zap.NewDevelopment()
	if opts.Environment == "production" {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatal(err)
	}

	if err := startHTTP(db, logger, opts.GoDiscoursePort); err != nil {
		log.Panicln(err)
	}
}
