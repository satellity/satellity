package models

import (
	"context"
	"fmt"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

const (
	GIST_GENRE_DEFAULT    = "DEFAULT"
	GIST_GENRE_RELEASE    = "RELEASE"
	GIST_GENRE_UPDATE     = "UPDATE"
	GIST_GENRE_NEWS       = "NEWS"
	GIST_GENRE_NEWSLETTER = "NEWSLETTER"
	GIST_GENRE_RESEARCH   = "RESEARCH"
)

type Gist struct {
	GistID    string
	Identity  string
	Author    string
	Title     string
	SourceID  string
	Genre     string
	Cardinal  bool
	Link      string
	Body      string
	PublishAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time

	Source *Source
}

var gistColumns = []string{"gist_id", "identity", "author", "title", "source_id", "genre", "cardinal", "link", "body", "publish_at", "created_at", "updated_at"}

func gistFromRow(row durable.Row) (*Gist, error) {
	var g Gist
	err := row.Scan(&g.GistID, &g.Identity, &g.Author, &g.Title, &g.SourceID, &g.Genre, &g.Cardinal, &g.Link, &g.Body, &g.PublishAt, &g.CreatedAt, &g.UpdatedAt)
	return &g, err
}

func (g *Gist) values() []any {
	return []any{g.GistID, g.Identity, g.Author, g.Title, g.SourceID, g.Genre, g.Cardinal, g.Link, g.Body, g.PublishAt, g.CreatedAt, g.UpdatedAt}
}

func CreateGist(ctx context.Context, identity, author, title, genre string, cardinal bool, link, body string, publishedAt time.Time, source *Source) (*Gist, error) {
	identity = strings.TrimSpace(identity)
	title = strings.TrimSpace(title)
	id := generateUniqueID(identity)

	t := time.Now()
	gist := &Gist{
		GistID:    id,
		Identity:  identity,
		Author:    author,
		Title:     title,
		SourceID:  source.SourceID,
		Genre:     genre,
		Cardinal:  cardinal,
		Link:      link,
		Body:      body,
		PublishAt: publishedAt,
		CreatedAt: t,
		UpdatedAt: t,
	}

	rows := [][]interface{}{
		gist.values(),
	}

	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		old, err := findGist(ctx, tx, id)
		if err != nil {
			return err
		} else if old != nil {
			gist = old
			return nil
		}
		old, err = findGistByLink(ctx, tx, gist.Link)
		if err != nil {
			return err
		} else if old != nil {
			gist = old
			return nil
		}
		_, err = tx.CopyFrom(ctx, pgx.Identifier{"gists"}, gistColumns, pgx.CopyFromRows(rows))
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return gist, nil
}

func (g *Gist) Update(ctx context.Context, title, genre string, cardinal bool) error {
	title = strings.TrimSpace(title)
	if title != "" {
		g.Title = title
	}
	if genre != "" {
		g.Genre = genre
	}
	g.Cardinal = cardinal

	cols, posits := durable.PrepareColumnsAndExpressions([]string{"title", "genre", "cardinal", "updated_at"}, 1)
	values := []interface{}{g.GistID, g.Title, g.Genre, g.Cardinal, time.Now()}
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf("UPDATE gists SET (%s)=(%s) WHERE gist_id=$1", cols, posits)
		_, err := tx.Exec(ctx, query, values...)
		return err
	})
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}

func ReadAllGists(ctx context.Context, offset time.Time) ([]*Gist, error) {
	if offset.IsZero() {
		offset = time.Now()
	}
	var gists []*Gist
	query := fmt.Sprintf("SELECT %s FROM gists WHERE publish_at<=$1 ORDER BY publish_at DESC LIMIT $2", strings.Join(gistColumns, ","))
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		gists, err = readGists(ctx, tx, query, offset, 512)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return gists, nil
}

func ReadGists(ctx context.Context, offset time.Time, limit int64) ([]*Gist, error) {
	if limit <= 0 || limit > 128 {
		limit = 128
	}
	if offset.IsZero() {
		offset = time.Now()
	}
	var gists []*Gist
	query := fmt.Sprintf("SELECT %s FROM gists WHERE cardinal=true AND publish_at<=$1 ORDER BY cardinal,publish_at DESC LIMIT $2", strings.Join(gistColumns, ","))
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		gists, err = readGists(ctx, tx, query, offset, limit)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return gists, nil
}

func readGists(ctx context.Context, tx pgx.Tx, query string, args ...interface{}) ([]*Gist, error) {
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gists []*Gist
	for rows.Next() {
		gist, err := gistFromRow(rows)
		if err != nil {
			return nil, err
		}
		gists = append(gists, gist)
	}
	sources, err := readSourceSet(ctx, tx, gists)
	if err != nil {
		return nil, err
	}
	for i, g := range gists {
		gists[i].Source = sources[g.SourceID]
	}
	return gists, nil
}

func ReadGist(ctx context.Context, id string) (*Gist, error) {
	var gist *Gist
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		gist, err = findGist(ctx, tx, id)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return gist, nil
}

func findGist(ctx context.Context, tx pgx.Tx, id string) (*Gist, error) {
	if uuid.FromStringOrNil(id).String() != id {
		return nil, nil
	}

	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM gists WHERE gist_id=$1", strings.Join(gistColumns, ",")), id)
	g, err := gistFromRow(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	g.Source, err = findSource(ctx, tx, g.SourceID)
	return g, err
}

func findGistByLink(ctx context.Context, tx pgx.Tx, link string) (*Gist, error) {
	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM gists WHERE link=$1", strings.Join(gistColumns, ",")), link)
	g, err := gistFromRow(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	g.Source, err = findSource(ctx, tx, g.SourceID)
	return g, err
}

func (g *Gist) Delete(ctx context.Context) error {
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "DELETE FROM gists WHERE gist_id=$1", g.GistID)
		return err
	})
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}
