package main

import (
	"context"
	"godiscourse/internal/durable"
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"
	_ "github.com/lib/pq"
)

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
	_, err := db.Exec("ALTER TABLE topics ALTER COLUMN short_id SET NOT NULL")
	if err != nil {
		log.Panicln(err)
	}
	_, err = db.Exec("CREATE UNIQUE INDEX ON topics (short_id);")
	if err != nil {
		log.Panicln(err)
	}
}

type Options struct {
	GoDiscourseURL  string `long:"url" env:"GODISCOURSE_URL" default:"http://localhost" required:"true"`
	GoDiscoursePort string `long:"port" env:"GODISCOURSE_PORT" default:"4000" requred:"true"`
	DbUser          string `long:"dbuser" env:"DB_USER" requred:"true"`
	DbPassword      string `long:"dbpassword" env:"DB_PASSWORD"`
	DbHost          string `long:"dbhost" env:"DB_HOST" default:"localhost"`
	DbPort          string `long:"dbport" env:"DB_PORT" default:"5432"`
	DbName          string `long:"dbname" env:"DB_NAME" default:"godiscourse_dev"`
	Environment     string `long:"environment" env:"ENV" default:"development"`
}
