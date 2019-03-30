package main

import (
	"context"
	"godiscourse/internal/durable"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	db := durable.OpenDatabaseClient(context.Background())
	_, err := db.Exec("ALTER TABLE topics ALTER COLUMN short_id SET NOT NULL")
	if err != nil {
		log.Panicln(err)
	}
	_, err = db.Exec("CREATE UNIQUE INDEX ON topics (short_id);")
	if err != nil {
		log.Panicln(err)
	}
}
