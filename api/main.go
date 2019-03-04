package main

import (
	"context"
	"log"

	"github.com/godiscourse/godiscourse/api/durable"
)

func main() {
	db := durable.OpenDatabaseClient(context.Background())
	defer db.Close()

	if err := startHTTP(db); err != nil {
		log.Panicln(err)
	}
}
