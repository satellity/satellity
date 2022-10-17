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
	log.Println("Feed service started at:", time.Now())
	ctx := session.WithDatabase(context.Background(), db)
	for {
		err := loopSources(ctx)
		if err != nil {
			log.Printf("services.loopSources error %s \n", err)
			time.Sleep(3 * time.Second)
			continue
		}
		time.Sleep(15 * time.Minute)
	}
}

func loopSources(ctx context.Context) error {
	sources, err := models.ReadSources(ctx)
	if err != nil {
		return err
	}
	for _, s := range sources {
		err = feeds.FetchCommon(ctx, s)
		if err != nil {
			log.Printf("loopSources link %s error %s \n", s.Link, err)
			time.Sleep(time.Second)
			s.Update(ctx, "", "", "", s.Wreck+1, time.Time{}, time.Time{})
			continue
		}
	}
	return nil
}
