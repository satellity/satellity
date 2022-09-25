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
			log.Printf("services.loopSources error %s \n", err)
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
		var err error
		switch s.Locality {
		case "github":
			err = feeds.FetchGithub(ctx, s)
		case "medium":
			err = feeds.FetchMedium(ctx, s)
		case "mirror":
			err = feeds.FetchMirror(ctx, s)
		case "substack":
			err = feeds.FetchSubStack(ctx, s)
		}
		if err != nil {
			log.Printf("loopSources link %s error %s \n", s.Link, err)
			time.Sleep(time.Second)
			continue
		}
	}
	return nil
}
