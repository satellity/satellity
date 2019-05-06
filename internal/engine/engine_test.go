package engine

import (
	"context"
	"godiscourse/internal/configs"
	"godiscourse/internal/durable"
	"godiscourse/internal/models"
	"log"
)

const (
	testEnvironment = "test"
	testDatabase    = "godiscourse_test"
)

func teardownTestContext(db *durable.Database) {
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
}

func setupTestContext() *durable.Database {
	opts := configs.DefaultOptions()
	if opts.Environment != testEnvironment {
		log.Panicln(opts.Environment)
	}
	db := durable.OpenDatabaseClient(context.Background(), &durable.ConnectionInfo{
		User:     opts.DbUser,
		Password: opts.DbPassword,
		Host:     opts.DbHost,
		Port:     opts.DbPort,
		Name:     opts.DbName,
	})
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
	return durable.WrapDatabase(db)
}
