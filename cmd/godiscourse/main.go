package main

import (
	"context"
	"database/sql"
	"fmt"
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

type Options struct {
	GoDiscourseURL  string `long:"url" env:"GODISCOURSE_URL" default:"http://localhost" required:"true"`
	GoDiscoursePort string `long:"port" env:"GODISCOURSE_PORT" default:"4000" requred:"true"`
	DbUser          string `long:"dbuser" env:"DB_USER" requred:"true"`
	DbPassword      string `long:"dbpassword" env:"DB_PASSWORD" required:"true"`
	DbHost          string `long:"dbhost" env:"DB_HOST" default:"localhost"`
	DbPort          string `long:"dbport" env:"DB_PORT" default:"5432"`
	DbName          string `long:"dbname" env:"DB_NAME" default:"godiscourse_dev"`
	Environment     string `long:"environment" env:"ENV" default:"development"`
}

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
	var opts Options
	p := flags.NewParser(&opts, flags.Default)

	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
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
