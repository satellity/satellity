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
	"path"
	"path/filepath"
	"strings"

	"github.com/dimfeld/httptreemux"
	"github.com/gorilla/handlers"
	flags "github.com/jessevdk/go-flags"
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

	return http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
}

func main() {
	var options struct {
		Dir         string `short:"d" long:"dir" env:"GO_DIR" description:"The config file directory"`
		Environment string `short:"e" long:"environment" env:"GO_ENV" default:"development"`
	}
	p := flags.NewParser(&options, flags.Default)
	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	if options.Dir == "" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Panicln(err)
		}
		back := ".."
		if strings.Contains(dir, "cmd") {
			back = "../.."
		}
		options.Dir = path.Join(dir, back, "internal/configs")
	}

	if err := configs.Init(options.Dir, options.Environment); err != nil {
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
		log.Fatal(err)
	}

	if err := startHTTP(db, logger, config.HTTP.Port); err != nil {
		log.Panicln(err)
	}
}
