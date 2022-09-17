package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"satellity/internal/configs"
	"satellity/internal/durable"
	"satellity/internal/middlewares"
	"satellity/internal/routes"

	"github.com/dimfeld/httptreemux"
	"github.com/gorilla/handlers"
	"github.com/jackc/pgx/v4/pgxpool"
	flags "github.com/jessevdk/go-flags"
	"github.com/unrolled/render"
	"go.uber.org/zap"
)

func startHTTP(db *pgxpool.Pool, logger *zap.Logger, port string) error {
	database := durable.WrapDatabase(db)
	router := httptreemux.New()
	routes.RegisterRoutes(router)

	handler := middlewares.Authenticate(router)
	handler = middlewares.Constraint(handler)
	handler = middlewares.Context(handler, database, render.New())
	handler = middlewares.State(handler)
	handler = middlewares.Logger(handler, durable.NewLogger(logger))
	handler = handlers.ProxyHeaders(handler)

	log.Printf("HTTP server running at: http://localhost:%s", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
}

func main() {
	var options struct {
		Environment string `short:"e" long:"environment" default:"development"`
	}
	p := flags.NewParser(&options, flags.Default)
	if _, err := p.Parse(); err != nil {
		log.Panicln(err)
	}

	if err := configs.Init(options.Environment); err != nil {
		log.Panicln(err)
	}

	config := configs.AppConfig
	db := durable.OpenDatabaseClient(context.Background(), &durable.ConnectionInfo{
		User:     config.Database.User,
		Password: config.Database.Password,
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		Name:     config.Database.Name,
	})
	defer db.Close()

	logger, err := zap.NewDevelopment()
	if config.Environment == "production" {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Panicln(err)
	}

	if err := startHTTP(db, logger, config.HTTP.Port); err != nil {
		log.Panicln(err)
	}
}
