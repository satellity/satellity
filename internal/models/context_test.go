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
	dropTopicUsersDDL = `DROP TABLE IF EXISTS topic_users;`
	dropTopicsDDL     = `DROP TABLE IF EXISTS topics;`
	dropCommentsDDL   = `DROP TABLE IF EXISTS comments;`
	dropStatisticsDDL = `DROP TABLE IF EXISTS statistics;`
)

func teardownTestContext(mctx *Context) {
	tables := []string{
		dropStatisticsDDL,
		dropCommentsDDL,
		dropTopicUsersDDL,
		dropTopicsDDL,
		dropCategoriesDDL,
		dropSessionsDDL,
		dropUsersDDL,
	}
	db := mctx.database
	for _, q := range tables {
		if _, err := db.Exec(q); err != nil {
			log.Panicln(err)
		}
	}
}

func setupTestContext() *Context {
	opts := configs.DefaultOptions()
	if opts.Environment != testEnvironment || opts.DbName != testDatabase {
		log.Panicln(opts.Environment, opts.DbName)
	}
	db := durable.OpenDatabaseClient(context.Background(), &durable.ConnectionInfo{
		User:     opts.DbUser,
		Password: opts.DbPassword,
		Host:     opts.DbHost,
		Port:     opts.DbPort,
		Name:     opts.DbName,
	})
	tables := []string{
		usersDDL,
		sessionsDDL,
		categoriesDDL,
		topicsDDL,
		topicUsersDDL,
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
