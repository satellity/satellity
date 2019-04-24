package main

import (
	"context"
	"godiscourse/internal/durable"
	"log"
	"os"
	"time"

	flags "github.com/jessevdk/go-flags"
	_ "github.com/lib/pq"
)

func main() {
	db := setupContext()
	offset := time.Now()
	limit := int64(50)
	for {
		count, last, err := MigrateTopics(db, offset, limit)
		if err != nil {
			log.Panicln(err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		offset = last
		if count < limit {
			break
		}
	}
}

func setupContext() *durable.Database {
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
	_, err := db.Exec("ALTER TABLE topics ADD COLUMN IF NOT EXISTS short_id VARCHAR(255)")
	if err != nil {
		log.Panicln(err)
	}
	return durable.WrapDatabase(db)
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
