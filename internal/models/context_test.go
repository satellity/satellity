package models

import (
	"context"
	"godiscourse/internal/configs"
	"godiscourse/internal/durable"
	"log"

	_ "github.com/lib/pq"
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

func teardownTestContext(context *Context) {
	tables := []string{
		dropStatisticsDDL,
		dropCommentsDDL,
		dropTopicsDDL,
		dropCategoriesDDL,
		dropSessionsDDL,
		dropUsersDDL,
	}
	db := context.database
	for _, q := range tables {
		if _, err := db.Exec(q); err != nil {
			log.Panicln(err)
		}
	}
}

func setupTestContext() *Context {
	if configs.Environment != testEnvironment {
		log.Panicln(configs.Environment)
	}
	//TODO: change db dependency in favor of in-memory structures
	db := durable.OpenDatabaseClient(context.Background(), &durable.ConnectionInfo{
		User:     "test",
		Password: "test",
		Host:     "localhost",
		Port:     "5432",
		Name:     "godicourse_test",
	})
	tables := []string{
		usersDDL,
		sessionsDDL,
		categoriesDDL,
		topicsDDL,
		commentsDDL,
		statisticsDDL,
	}
	for _, q := range tables {
		if _, err := db.Exec(q); err != nil {
			log.Panicln(err)
		}
	}
	database := durable.WrapDatabase(db)
	return WrapContext(context.Background(), database)
}
