package models

import (
	"context"
	"log"
	"satellity/internal/configs"
	"satellity/internal/durable"
)

const (
	testEnvironment = "test"
	testDatabase    = "satellity_test"
)

const (
	dropUsersDDL             = `DROP TABLE IF EXISTS users;`
	dropSessionsDDL          = `DROP TABLE IF EXISTS sessions;`
	dropEmailVerificationDDL = `DROP TABLE IF EXISTS email_verifications;`
	dropCategoriesDDL        = `DROP TABLE IF EXISTS categories;`
	dropTopicUsersDDL        = `DROP TABLE IF EXISTS topic_users;`
	dropTopicsDDL            = `DROP TABLE IF EXISTS topics;`
	dropCommentsDDL          = `DROP TABLE IF EXISTS comments;`
	dropStatisticsDDL        = `DROP TABLE IF EXISTS statistics;`
)

func teardownTestContext(mctx *Context) {
	tables := []string{
		dropStatisticsDDL,
		dropCommentsDDL,
		dropTopicUsersDDL,
		dropTopicsDDL,
		dropCategoriesDDL,
		dropEmailVerificationDDL,
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
	if err := configs.Init("../configs/config.yaml", testEnvironment); err != nil {
		log.Panicln(err)
	}
	config := configs.AppConfig
	if config.Environment != testEnvironment || config.Database.Name != testDatabase {
		log.Panicln(config.Environment, config.Database.Name)
	}
	db := durable.OpenDatabaseClient(context.Background(), &durable.ConnectionInfo{
		User:     config.Database.User,
		Password: config.Database.Password,
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		Name:     config.Database.Name,
	})
	tables := []string{
		usersDDL,
		sessionsDDL,
		emailVerificationDDL,
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
