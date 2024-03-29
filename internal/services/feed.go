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
		var err error
		switch s.Locality {
		case "github":
			err = feeds.FetchGithub(ctx, s)
		case "medium":
			err = feeds.FetchMedium(ctx, s)
		case "mirror":
			err = feeds.FetchMirror(ctx, s)
		case "messari",
			"substack",
			"techcrunch",
			"decrypt",
			"crunchbase",
			"trustnodes",
			"common":
			err = feeds.FetchCommon(ctx, s)
		}
		if err != nil {
			log.Printf("loopSources link %s error %s \n", s.Link, err)
			time.Sleep(time.Second)
			s.Update(ctx, "", "", "", s.Wreck+1, time.Time{}, time.Time{})
			continue
		}
	}
	return nil
}
