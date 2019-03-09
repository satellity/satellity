package durable

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/godiscourse/godiscourse/api/config"
	_ "github.com/lib/pq" //
)

type Database struct {
	db *sql.DB
}

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

func WrapDatabase(db *sql.DB) *Database {
	return &Database{db: db}
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Exec(args...)
}

func (d *Database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.ExecContext(ctx, args...)
}

func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Query(args...)
}

func (d *Database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryContext(ctx, args...)
}

func (d *Database) QueryRow(query string, args ...interface{}) (*sql.Row, error) {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRow(args...), nil
}

func (d *Database) QueryRowContext(ctx context.Context, query string, args ...interface{}) (*sql.Row, error) {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRowContext(ctx, args...), nil
}

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

func PrepareColumnsWithValues(columns []string) (string, string) {
	if len(columns) < 1 {
		return "", ""
	}
	cols, params := bytes.Buffer{}, bytes.Buffer{}
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

type Row interface {
	Scan(dest ...interface{}) error
}
