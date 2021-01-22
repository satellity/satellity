package models

import (
	"context"
	"io/ioutil"
	"log"
	"satellity/internal/configs"
	"satellity/internal/durable"
	"satellity/internal/session"
)

const (
	testEnvironment = "test"
	testDatabase    = "satellity_test"

	dropCategoriesDDL        = `DROP TABLE IF EXISTS categories;`
	dropCommentsDDL          = `DROP TABLE IF EXISTS comments;`
	dropEmailVerificationDDL = `DROP TABLE IF EXISTS email_verifications;`
	dropSessionsDDL          = `DROP TABLE IF EXISTS sessions;`
	dropStatisticsDDL        = `DROP TABLE IF EXISTS statistics;`
	dropProductsDDL          = `DROP TABLE IF EXISTS products;`
	dropTopicsDDL            = `DROP TABLE IF EXISTS topics;`
	dropTopicUsersDDL        = `DROP TABLE IF EXISTS topic_users;`
	dropUsersDDL             = `DROP TABLE IF EXISTS users;`
)

func teardownTestContext(ctx context.Context) {
	tables := []string{
		dropStatisticsDDL,
		dropProductsDDL,
		dropCommentsDDL,
		dropTopicUsersDDL,
		dropTopicsDDL,
		dropCategoriesDDL,
		dropEmailVerificationDDL,
		dropSessionsDDL,
		dropUsersDDL,
	}
	db := session.Database(ctx)
	for _, q := range tables {
		if _, err := db.Exec(ctx, q); err != nil {
			log.Panicln(err)
		}
	}
}

func setupTestContext() context.Context {
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
	data, err := ioutil.ReadFile("./schema.sql")
	if err != nil {
		log.Panicln(err)
	}
	if _, err := db.Exec(context.Background(), string(data)); err != nil {
		log.Panicln(err)
	}
	database := durable.WrapDatabase(db)
	return session.WithDatabase(context.Background(), database)
}
