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

type Origin struct {
	OriginID   string
	Name       string
	WebsiteURL string
	LogoURL    string
	Locality   string
	Wreck      int
	Position   int64
	CreatedAt  time.Time
	UpdateAt   time.Time
}

var originColumns = []string{"origin_id", "name", "website_url", "logo_url", "locality", "wreck", "position", "created_at", "updated_at"}

func (r *Origin) values() []any {
	return []any{r.OriginID, r.Name, r.WebsiteURL, r.LogoURL, r.Locality, r.Wreck, r.Position, r.CreatedAt, r.UpdateAt}
}

func originFromRows(row durable.Row) (*Origin, error) {
	var o Origin
	err := row.Scan(&o.OriginID, &o.Name, &o.WebsiteURL, &o.LogoURL, &o.Locality, &o.Wreck, &o.Position, &o.CreatedAt, &o.UpdateAt)
	return &o, err
}

func CreateOrigin(ctx context.Context, name, link, logo, locality string) (*Origin, error) {
	_, err := url.Parse(link)
	if err != nil {
		return nil, session.BadDataErrorWithFieldAndData(ctx, "link", "invalid", link)
	}

	name = strings.TrimSpace(name)
	link = strings.TrimSpace(link)
	logo = strings.TrimSpace(logo)
	locality = strings.TrimSpace(locality)

	id := generateUniqueID(link)
	t := time.Now()
	origin := &Origin{
		OriginID:   id,
		Name:       name,
		WebsiteURL: link,
		LogoURL:    logo,
		Locality:   locality,
		Wreck:      0,
		Position:   10001,
		CreatedAt:  t,
		UpdateAt:   t,
	}

	rows := [][]interface{}{
		origin.values(),
	}
	err = session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		old, err := findOrigin(ctx, tx, id)
		if err != nil {
			return err
		} else if old != nil {
			origin = old
			return nil
		}
		_, err = tx.CopyFrom(ctx, pgx.Identifier{"origins"}, originColumns, pgx.CopyFromRows(rows))
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return origin, nil
}

func ReadOrigin(ctx context.Context, id string) (*Origin, error) {
	var origin *Origin
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		origin, err = findOrigin(ctx, tx, id)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return origin, nil
}

func findOrigin(ctx context.Context, tx pgx.Tx, id string) (*Origin, error) {
	if uuid.FromStringOrNil(id).String() != id {
		return nil, nil
	}

	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM origins WHERE origin_id=$1", strings.Join(originColumns, ",")), id)
	o, err := originFromRows(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return o, err
}
