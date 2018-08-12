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
	drop_users_DDL = `DROP TABLE IF EXISTS users;`
)

func teardownTestContext(ctx context.Context) {
	tables := []string{
		drop_users_DDL,
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
		users_DDL,
	}
	for _, q := range tables {
		db.Exec(q)
	}
	return session.WithDatabase(context.Background(), db)
}
