package durable

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" //
)

type ConnectionInfo struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

// Database is wrapped struct of *sql.DB
type Database struct {
	db *sql.DB
}

// OpenDatabaseClient generate a database client
func OpenDatabaseClient(ctx context.Context, c *ConnectionInfo) *sql.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Name)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	if err := db.Ping(); err != nil {
		log.Fatal(fmt.Errorf("\nFail to connect the database.\nPlease make sure the connection info is valid %#v", c))
		return nil
	}
	return db
}

// WrapDatabase create a *Database
func WrapDatabase(db *sql.DB) *Database {
	return &Database{db: db}
}

// Close the *sql.DB
func (d *Database) Close() error {
	return d.db.Close()
}

// Exec executes a prepared statement
func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Exec(args...)
}

// ExecContext executes a prepared statement with a context
func (d *Database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.ExecContext(ctx, args...)
}

// Query executes a prepared query statement with the given arguments
func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Query(args...)
}

// QueryContext executes a prepared query statement with the given arguments
func (d *Database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryContext(ctx, args...)
}

// QueryRow executes a prepared query statement with the given arguments.
func (d *Database) QueryRow(query string, args ...interface{}) (*sql.Row, error) {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRow(args...), nil
}

// QueryRowContext executes a prepared query statement with the given arguments.
func (d *Database) QueryRowContext(ctx context.Context, query string, args ...interface{}) (*sql.Row, error) {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRowContext(ctx, args...), nil
}

// RunInTransaction run a query in the transaction
func (d *Database) RunInTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			_ = tx.Rollback()
			panic(err)
		}
	}()
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

// PrepareColumnsWithValues prepare columns and placeholder
func PrepareColumnsWithValues(columns []string) (string, string) {
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
