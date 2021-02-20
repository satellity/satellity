package durable

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// ConnectionInfo is the info of the postgres
type ConnectionInfo struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

// Database is wrapped struct of *pgx.Conn
type Database struct {
	db *pgxpool.Pool
}

// OpenDatabaseClient generate a database client
func OpenDatabaseClient(ctx context.Context, c *ConnectionInfo) *pgxpool.Pool {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Name)
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Panicln(err)
	}
	config.MinConns = 1
	config.MaxConns = 128
	dbpool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Panicln(err)
	}
	if err := dbpool.Ping(ctx); err != nil {
		log.Fatal(fmt.Errorf("\nFail to connect the database.\nPlease make sure the connection info is valid %#v", c))
		return nil
	}
	return dbpool
}

// WrapDatabase create a *Database
func WrapDatabase(db *pgxpool.Pool) *Database {
	return &Database{db: db}
}

// Close the *pgx.Conn
func (d *Database) Close() {
	d.db.Close()
}

// Exec executes a prepared statement
func (d *Database) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return d.db.Exec(ctx, query, args...)
}

// Query executes a prepared query statement with the given arguments
func (d *Database) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return d.db.Query(ctx, query, args...)
}

// QueryRowContext executes a prepared query statement with the given arguments.
func (d *Database) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return d.db.QueryRow(ctx, query, args...)
}

// RunInTransaction run a query in the transaction
func (d *Database) RunInTransaction(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := d.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// PrepareColumnsWithParams prepare columns and placeholders
func PrepareColumnsWithParams(columns []string) (string, string) {
	if len(columns) < 1 {
		return "", ""
	}
	var cols, params bytes.Buffer
	for i, column := range columns {
		if i > 0 {
			cols.WriteString(",")
			params.WriteString(",")
		}
		cols.WriteString(column)
		params.WriteString(fmt.Sprintf("$%d", i+1))
	}
	return cols.String(), params.String()
}

// Row is a interface
type Row interface {
	Scan(dest ...interface{}) error
}
