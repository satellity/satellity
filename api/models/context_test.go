package models

import (
	"context"
	"log"

	"github.com/godiscourse/godiscourse/api/config"
	"github.com/godiscourse/godiscourse/api/durable"
	"github.com/godiscourse/godiscourse/api/session"
)

const (
	testEnvironment = "test"
	testDatabase    = "godiscourse_test"
)

const (
	dropUsersDDL      = `DROP TABLE IF EXISTS users;`
	dropSessionsDDL   = `DROP TABLE IF EXISTS sessions;`
	dropCategoriesDDL = `DROP TABLE IF EXISTS categories;`
	dropTopicsDDL     = `DROP TABLE IF EXISTS topics;`
	dropCommentsDDL   = `DROP TABLE IF EXISTS comments;`
	dropStatisticsDDL = `DROP TABLE IF EXISTS statistics;`
)

func teardownTestContext(ctx context.Context) {
	tables := []string{
		dropStatisticsDDL,
		dropCommentsDDL,
		dropTopicsDDL,
		dropCategoriesDDL,
		dropSessionsDDL,
		dropUsersDDL,
	}
	for _, q := range tables {
		session.Database(ctx).Exec(q)
	}
}

func setupTestContext() context.Context {
	if config.Environment != testEnvironment || config.DatabaseName != testDatabase {
		log.Panicln(config.Environment, config.DatabaseName)
	}
	db := durable.OpenDatabaseClient(context.Background())
	tables := []string{
		usersDDL,
		sessionsDDL,
		categoriesDDL,
		topicsDDL,
		commentsDDL,
		statisticsDDL,
	}
	for _, q := range tables {
		db.Exec(q)
	}
	return session.WithDatabase(context.Background(), db)
}
