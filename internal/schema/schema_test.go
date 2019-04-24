package schema

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

func teardownTestContext(db *durable.Database) {
	tables := []string{
		dropStatisticsDDL,
		dropCommentsDDL,
		dropTopicsDDL,
		dropCategoriesDDL,
		dropSessionsDDL,
		dropUsersDDL,
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
	return durable.WrapDatabase(db)
}
