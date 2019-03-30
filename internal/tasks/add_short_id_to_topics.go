package main

import (
	"context"
	"godiscourse/internal/durable"
	"godiscourse/internal/models"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	ctx := setupContext()
	offset := time.Now()
	limit := int64(50)
	for {
		count, last, err := models.MigrateTopics(ctx, offset, limit)
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

func setupContext() *models.Context {
	db := durable.OpenDatabaseClient(context.Background())
	_, err := db.Exec("ALTER TABLE topics ADD COLUMN IF NOT EXISTS short_id VARCHAR(255)")
	if err != nil {
		log.Panicln(err)
	}
	database := durable.WrapDatabase(db)
	return models.WrapContext(context.Background(), database)
}
