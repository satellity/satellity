package engine

import (
	"context"
	"godiscourse/internal/configs"
	"godiscourse/internal/durable"
	"godiscourse/internal/models"
	"log"
)

func SetupTestContext() (*durable.Database, func()) {
	opts := configs.DefaultOptions()

	db := durable.OpenDatabaseClient(context.Background(), &durable.ConnectionInfo{
		User:     opts.DbUser,
		Password: opts.DbPassword,
		Host:     opts.DbHost,
		Port:     opts.DbPort,
		Name:     opts.DbName,
	})

	if _, err := db.Exec("CREATE DATABASE godiscourse_test"); err != nil {
		log.Panicln(err)
	}

	tables := []string{
		models.UsersDDL,
		models.SessionsDDL,
		models.CategoriesDDL,
		models.TopicsDDL,
		models.CommentsDDL,
		models.StatisticsDDL,
	}
	for _, q := range tables {
		if _, err := db.Exec(q); err != nil {
			log.Panicln(err)
		}
	}

	teardown := func() {
		tables := []string{
			models.DropStatisticsDDL,
			models.DropCommentsDDL,
			models.DropTopicsDDL,
			models.DropCategoriesDDL,
			models.DropSessionsDDL,
			models.DropUsersDDL,
		}

		for _, q := range tables {
			if _, err := db.Exec(q); err != nil {
				log.Panicln(err)
			}
		}

		if _, err := db.Exec("DROP DATABASE godiscourse_test"); err != nil {
			log.Panicln(err)
		}
	}

	return durable.WrapDatabase(db), teardown
}
