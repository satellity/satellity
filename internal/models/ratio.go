package models

import (
	"context"
	"fmt"
	"log"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

const (
	TopLongShortAccountRatio    = "TOP_LONG_SHORT_ACCOUNT_RATIO"
	TopLongShortPositionRatio   = "TOP_LONG_SHORT_POSITION_RATIO"
	GlobalLongShortAccountRatio = "GLOBAL_LONG_SHORT_ACCOUNT_RATIO"

	RatioPeriod5M = "5M"
)

type Ratio struct {
	RatioID        string
	Category       string
	Symbol         string
	Period         string
	LongAccount    string
	ShortAccount   string
	LongShortRatio string
	TimestampAt    int64
}

var ratioColumns = []string{"ratio_id", "category", "symbol", "period", "long_account", "short_account", "long_short_ratio", "timestamp_at"}

func ratioFromRows(row durable.Row) (*Ratio, error) {
	var r Ratio
	err := row.Scan(&r.RatioID, &r.Category, &r.Symbol, &r.Period, &r.LongAccount, &r.ShortAccount, &r.LongShortRatio, &r.TimestampAt)
	return &r, err
}

func (r *Ratio) values() []any {
	return []any{r.RatioID, r.Category, r.Symbol, r.Period, r.LongAccount, r.ShortAccount, r.LongShortRatio, r.TimestampAt}
}

func UpsertRatio(ctx context.Context, category, symbol, period, long, short, lsRatio string, at int64) (*Ratio, error) {
	ratio := &Ratio{
		RatioID:        generateUniqueID(category, symbol, period, fmt.Sprint(at)),
		Category:       category,
		Symbol:         symbol,
		Period:         period,
		LongAccount:    long,
		ShortAccount:   short,
		LongShortRatio: lsRatio,
		TimestampAt:    at,
	}

	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		old, err := findRatio(ctx, tx, ratio.RatioID)
		if err != nil {
			return err
		} else if old != nil {
			ratio = old
			return nil
		}

		rows := [][]interface{}{
			ratio.values(),
		}
		_, err = tx.CopyFrom(ctx, pgx.Identifier{"ratios"}, ratioColumns, pgx.CopyFromRows(rows))
		return err
	})
	if err != nil {
		log.Println("session.TransactionError:", err)
		return nil, session.TransactionError(ctx, err)
	}
	return ratio, nil
}

func findRatio(ctx context.Context, tx pgx.Tx, id string) (*Ratio, error) {
	if uuid.FromStringOrNil(id).String() != id {
		return nil, nil
	}

	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM ratios WHERE ratio_id=$1", strings.Join(ratioColumns, ",")), id)
	r, err := ratioFromRows(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return r, err
}

func FillAssetWithRatios(ctx context.Context, assets []*Asset) ([]*Asset, error) {
	symbols := make([]string, len(assets))
	for i, a := range assets {
		symbols[i] = a.Contract
	}
	categories := []string{
		TopLongShortAccountRatio,
		TopLongShortPositionRatio,
		GlobalLongShortAccountRatio,
	}
	set := make(map[string][]*Ratio, 0)
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, fmt.Sprintf("SELECT %s FROM (SELECT *, rank() OVER (PARTITION BY category,symbol,period ORDER BY timestamp_at DESC) FROM ratios WHERE category=ANY($1) AND symbol=ANY($2) AND period=$3 AND timestamp_at>$4) filter WHERE RANK=1", strings.Join(ratioColumns, ",")), categories, symbols, RatioPeriod5M, time.Now().Add(-30*time.Minute).Unix()*1000)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			ratio, err := ratioFromRows(rows)
			if err != nil {
				return err
			}
			set[ratio.Symbol] = append(set[ratio.Symbol], ratio)
		}
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	for _, asset := range assets {
		asset.Ratios = set[asset.Contract]
	}
	return assets, nil
}
