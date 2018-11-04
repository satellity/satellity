package durable

import (
	"context"
	"fmt"

	"github.com/go-pg/pg"
	"github.com/godiscourse/godiscourse/config"
)

// OpenDatabaseClient generate a database client
func OpenDatabaseClient(ctx context.Context) *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", config.DatabaseHost, config.DatabasePort),
		User:     config.DatabaseUser,
		Password: config.DatabasePassword,
		Database: config.DatabaseName,
	})
	return db
}
