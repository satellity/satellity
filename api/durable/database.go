package durable

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/godiscourse/godiscourse/api/config"
	_ "github.com/lib/pq" //
)

// OpenDatabaseClient generate a database client
func OpenDatabaseClient(ctx context.Context) *sql.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.DatabaseUser, config.DatabasePassword, config.DatabaseHost, config.DatabasePort, config.DatabaseName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

type Row interface {
	Scan(dest ...interface{}) error
}
