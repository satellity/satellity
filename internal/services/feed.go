package services

import (
	"context"
	"log"
	"satellity/internal/durable"
	"satellity/internal/models"
	"satellity/internal/services/feeds"
	"satellity/internal/session"
	"time"
)

func Run(db *durable.Database) {
	ctx := session.WithDatabase(context.Background(), db)
	for {
		err := loopSources(ctx)
		if err != nil {
			log.Printf("services.loopSources error %s", err)
			time.Sleep(3 * time.Second)
			continue
		}
		time.Sleep(time.Hour)
	}
}

func loopSources(ctx context.Context) error {
	sources, err := models.ReadSources(ctx)
	if err != nil {
		return err
	}
	for _, s := range sources {
		log.Println("fetch", s.Link)
		err := feeds.Release(ctx, s.Link)
		if err != nil {
			return err
		}
		break
	}
	return nil
}
