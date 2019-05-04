package main

import (
	"context"
	"database/sql"
	"fmt"
	"godiscourse/internal/category"
	"godiscourse/internal/comment"
	"godiscourse/internal/configs"
	"godiscourse/internal/controllers"
	"godiscourse/internal/controllers/admin"
	"godiscourse/internal/durable"
	"godiscourse/internal/middleware"
	"godiscourse/internal/topic"
	"godiscourse/internal/user"
	"log"
	"net/http"
	"os"

	"github.com/dimfeld/httptreemux"
	"github.com/gorilla/handlers"
	flags "github.com/jessevdk/go-flags"
	_ "github.com/lib/pq"
	"github.com/unrolled/render"
	"go.uber.org/zap"
)

func startHTTP(db *sql.DB, logger *zap.Logger, port string) error {
	database := durable.WrapDatabase(db)
	u := user.New(database)
	cat := category.New(database)
	t := topic.New(database, u, cat)
	c := comment.New(database, u, t)

	router := httptreemux.New()

	controllers.RegisterHanders(router)
	controllers.HealthCheck(router)
	controllers.RegisterUser(u, t, router)
	controllers.RegisterCategory(cat, t, router)
	controllers.RegisterTopic(t, router)
	controllers.RegisterComment(c, t, router)

	admin.RegisterAdminUser(u, router)
	admin.RegisterAdminCategory(cat, router)
	admin.RegisterAdminTopic(t, router)

	handler := middleware.Authenticate(u, router)
	handler = middleware.Constraint(handler)
	handler = middleware.Context(handler, render.New())
	handler = middleware.State(handler)
	handler = middleware.Logger(handler, durable.NewLogger(logger))
	handler = handlers.ProxyHeaders(handler)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
}

func main() {
	opts := configs.DefaultOptions()
	if configs.Environment == "production" {
		p := flags.NewParser(opts, flags.Default)
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
