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
	LogoURL   string
	Locality  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var sourceColumns = []string{"source_id", "author", "host", "link", "logo_url", "locality", "created_at", "updated_at"}

func (s *Source) values() []any {
	return []any{s.SourceID, s.Author, s.Host, s.Link, s.LogoURL, s.Locality, s.CreatedAt, s.UpdatedAt}
}

func sourceFromRows(row durable.Row) (*Source, error) {
	var s Source
	err := row.Scan(&s.SourceID, &s.Author, &s.Host, &s.Link, &s.LogoURL, &s.Locality, &s.CreatedAt, &s.UpdatedAt)
	return &s, err
}

func CreateSource(ctx context.Context, author, link, logo, locality string) (*Source, error) {
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

	id := generateUniqueID(link)

	t := time.Now()
	source := &Source{
		SourceID:  id,
		Author:    author,
		Host:      host,
		Link:      link,
		LogoURL:   logo,
		Locality:  locality,
		CreatedAt: t,
		UpdatedAt: t,
	}

	rows := [][]interface{}{
		source.values(),
	}
	err = session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		old, err := findSource(ctx, tx, id)
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

func (s *Source) Update(ctx context.Context, author, host, logo string, updated time.Time) error {
	author = strings.TrimSpace(author)
	host = strings.TrimSpace(host)
	logo = strings.TrimSpace(logo)
	if author != "" {
		s.Author = author
	}
	if host != "" {
		s.Host = host
	}
	if logo != "" {
		s.LogoURL = logo
	}

	cols, posits := durable.PrepareColumnsAndExpressions([]string{"author", "host", "logo_url", "updated_at"}, 1)
	values := []interface{}{s.SourceID, s.Author, s.Host, s.LogoURL, updated}
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
	rows, err := tx.Query(ctx, fmt.Sprintf("SELECT %s FROM sources ORDER BY updated_at LIMIT 3000", strings.Join(sourceColumns, ",")))
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

func readSourceSet(ctx context.Context, tx pgx.Tx, gists []*Gist) (map[string]*Source, error) {
	ids := make([]string, len(gists))
	for i, g := range gists {
		ids[i] = g.SourceID
	}
	set := make(map[string]*Source)
	if len(ids) < 1 {
		return set, nil
	}
	query := fmt.Sprintf("SELECT %s FROM sources WHERE source_id=ANY($1)", strings.Join(sourceColumns, ","))
	rows, err := tx.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		source, err := sourceFromRows(rows)
		if err != nil {
			return nil, err
		}
		set[source.SourceID] = source
	}
	return set, rows.Err()
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

func (s *Source) Delete(ctx context.Context) error {
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "DELETE FROM sources WHERE source_id=$1", s.SourceID)
		return err
	})
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}
