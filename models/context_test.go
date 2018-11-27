package models

import (
	"context"
	"log"

	"github.com/godiscourse/godiscourse/config"
	"github.com/godiscourse/godiscourse/durable"
	"github.com/godiscourse/godiscourse/session"
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
)

func teardownTestContext(ctx context.Context) {
	tables := []string{
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
	}
	for _, q := range tables {
		db.Exec(q)
	}
	return session.WithDatabase(context.Background(), db)
}
