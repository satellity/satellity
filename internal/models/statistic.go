package models

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

// SolidStatisticID is used to generate a solid id from name
const SolidStatisticID = "540cbd3c-f4eb-479c-bcd8-b5629af57267"

const (
	StatisticTypeUsers    = "users"
	StatisticTypeTopics   = "topics"
	StatisticTypeComments = "comments"
)

// Statistic is the body of statistic
type Statistic struct {
	StatisticID string
	Name        string
	Count       int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var statisticColumns = []string{"statistic_id", "name", "count", "created_at", "updated_at"}

func (s *Statistic) values() []interface{} {
	return []interface{}{s.StatisticID, s.Name, s.Count, s.CreatedAt, s.UpdatedAt}
}

func statisticFromRows(row durable.Row) (*Statistic, error) {
	var s Statistic
	err := row.Scan(&s.StatisticID, &s.Name, &s.Count, &s.CreatedAt, &s.UpdatedAt)
	return &s, err
}

func UpsertStatistic(ctx context.Context, name string) (*Statistic, error) {
	id, err := generateStatisticID(SolidStatisticID, name)
	if err != nil {
		return nil, session.ServerError(ctx, err)
	}
	switch name {
	case StatisticTypeUsers,
		StatisticTypeTopics,
		StatisticTypeComments:
	default:
		return nil, session.BadDataError(ctx)
	}

	var statistic *Statistic
	err = session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		s, err := findStatistic(ctx, tx, id)
		if err != nil {
			return err
		}
		var count int64
		switch name {
		case "users":
			count, err = usersCount(ctx, tx)
		case "topics":
			count, err = topicsCount(ctx, tx)
		case "comments":
			count, err = commentsCount(ctx, tx)
		}
		if err != nil {
			return err
		}
		t := time.Now()
		if s == nil {
			s = &Statistic{
				StatisticID: id,
				Name:        name,
				CreatedAt:   t,
			}
		}
		s.Count = count
		s.UpdatedAt = t
		cols, params := durable.PrepareColumnsAndExpressions(statisticColumns, 0)
		_, err = tx.Exec(ctx, fmt.Sprintf("INSERT INTO statistics(%s) VALUES (%s) ON CONFLICT (statistic_id) DO UPDATE SET (count,updated_at)=(EXCLUDED.count,EXCLUDED.updated_at)", cols, params), s.values()...)
		statistic = s
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return statistic, nil
}

func findStatistic(ctx context.Context, tx pgx.Tx, id string) (*Statistic, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}

	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM Statistics WHERE statistic_id=$1", strings.Join(statisticColumns, ",")), id)
	s, err := statisticFromRows(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func generateStatisticID(ID, name string) (string, error) {
	h := md5.New()
	io.WriteString(h, ID)
	io.WriteString(h, name)
	sum := h.Sum(nil)
	sum[6] = (sum[6] & 0x0f) | 0x30
	sum[8] = (sum[8] & 0x3f) | 0x80
	id, err := uuid.FromBytes(sum)
	return id.String(), err
}
