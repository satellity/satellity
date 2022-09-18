package models

import (
	"context"
	"fmt"
	"net/url"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

type Source struct {
	SourceID  string
	Author    string
	Host      string
	Link      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var sourceColumns = []string{"source_id", "author", "host", "link", "created_at", "updated_at"}

func (s *Source) values() []any {
	return []any{s.SourceID, s.Author, s.Host, s.Link, s.CreatedAt, s.UpdatedAt}
}

func sourceFromRows(row durable.Row) (*Source, error) {
	var s Source
	err := row.Scan(&s.SourceID, &s.Author, &s.Host, &s.Link, &s.CreatedAt, &s.UpdatedAt)
	return &s, err
}

func CreateSource(ctx context.Context, author, link string) (*Source, error) {
	author = strings.TrimSpace(author)
	link = strings.TrimSpace(link)

	if author == "" {
		return nil, session.BadDataErrorWithFieldAndData(ctx, "author", "invalid", author)
	}
	if link == "" {
		return nil, session.BadDataErrorWithFieldAndData(ctx, "link", "invalid", author)
	}

	uri, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	host := uri.Host
	host = strings.Replace(host, "www.", "", 0)

	t := time.Now()
	source := &Source{
		SourceID:  uuid.Must(uuid.NewV4()).String(),
		Author:    author,
		Link:      link,
		CreatedAt: t,
		UpdatedAt: t,
	}

	rows := [][]interface{}{
		source.values(),
	}
	err = session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		old, err := findSourceByLink(ctx, tx, link)
		if err != nil {
			return err
		} else if old != nil {
			source = old
			return nil
		}
		_, err = tx.CopyFrom(ctx, pgx.Identifier{"sources"}, sourceColumns, pgx.CopyFromRows(rows))
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return source, nil
}

func (s *Source) Update(ctx context.Context, author, host string) error {
	author = strings.TrimSpace(author)
	host = strings.TrimSpace(host)
	if author != "" {
		s.Author = author
	}
	if host != "" {
		s.Host = host
	}

	cols, posits := durable.PrepareColumnsAndExpressions([]string{"author", "host"}, 1)
	values := []interface{}{s.SourceID, s.Author, s.Host}
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf("UPDATE sources SET (%s)=(%s) WHERE source_id=$1", cols, posits)
		_, err := tx.Exec(ctx, query, values...)
		return err
	})
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}

func ReadSources(ctx context.Context) ([]*Source, error) {
	var sources []*Source
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		sources, err = readSources(ctx, tx)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return sources, nil
}

func readSources(ctx context.Context, tx pgx.Tx) ([]*Source, error) {
	rows, err := tx.Query(ctx, fmt.Sprintf("SELECT %s FROM sources LIMIT 3000", strings.Join(sourceColumns, ",")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []*Source
	for rows.Next() {
		source, err := sourceFromRows(rows)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}
	return sources, rows.Err()
}

func ReadSource(ctx context.Context, id string) (*Source, error) {
	var source *Source
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		source, err = findSource(ctx, tx, id)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return source, nil
}

func findSource(ctx context.Context, tx pgx.Tx, id string) (*Source, error) {
	if uuid.FromStringOrNil(id).String() != id {
		return nil, nil
	}

	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM sources WHERE source_id=$1", strings.Join(sourceColumns, ",")), id)
	s, err := sourceFromRows(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func findSourceByLink(ctx context.Context, tx pgx.Tx, link string) (*Source, error) {
	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM sources WHERE link=$1", strings.Join(sourceColumns, ",")), link)
	s, err := sourceFromRows(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return s, err
}
