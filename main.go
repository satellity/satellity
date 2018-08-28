package main

import (
	"context"
	"log"

	"github.com/godiscourse/godiscourse/durable"
)

func main() {
	db := durable.OpenDatabaseClient(context.Background())
	defer db.Close()

	if err := startHttp(db); err != nil {
		log.Panicln(err)
	}
}
