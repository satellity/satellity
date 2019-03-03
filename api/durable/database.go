package durable

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/godiscourse/godiscourse/api/config"
	_ "github.com/lib/pq"
)

// OpenDatabaseClient generate a database client
func OpenDatabaseClient(ctx context.Context) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password='%s' dbname=%s sslmode=disable", config.DatabaseHost, config.DatabasePort, config.DatabaseUser, config.DatabasePassword, config.DatabaseName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}
