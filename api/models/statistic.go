package models

import (
	"context"
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/godiscourse/godiscourse/api/session"
	"github.com/gofrs/uuid"
)

// SolidStatisticID is used to generate a solid id from name
const SolidStatisticID = "540cbd3c-f4eb-479c-bcd8-b5629af57267"

const statisticsDDL = `
CREATE TABLE IF NOT EXISTS statistics (
	statistic_id          VARCHAR(36) PRIMARY KEY,
	name                  VARCHAR(512) NOT NULL,
	count                 BIGINT NOT NULL DEFAULT 0,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
`

// Statistic is the body of statistic
type Statistic struct {
	StatisticID string    `sql:"statistic_id,pk"`
	Name        string    `sql:"name,notnull"`
	Count       int64     `sql:"count,notnull"`
	CreatedAt   time.Time `sql:"created_at"`
	UpdatedAt   time.Time `sql:"updated_at"`
}

var statisticColums = []string{"statistic_id", "name", "count", "created_at", "updated_at"}

func (s *Statistic) values() []interface{} {
	return []interface{}{s.StatisticID, s.Name, s.Count, s.CreatedAt, s.UpdatedAt}
}

func upsertStatistic(ctx context.Context, tx *sql.Tx, name string) (*Statistic, error) {
	id, _ := generateStatisticID(SolidStatisticID, name)
	s, err := findStatistic(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	var count int64
	switch name {
	case "users":
		count, err = usersCount(ctx, tx)
	case "topics":
		count, err = topicsCount(ctx, tx)
	case "comments":
		count, err = commentsCount(ctx, tx)
	default:
		return nil, session.BadDataError(ctx)
	}
	if err != nil {
		return nil, err
	}
	if s != nil {
		s.Count = count
		_, err := tx.ExecContext(ctx, fmt.Sprintf("UPDATE statistics SET count=$1 WHERE statistic_id=$2"), count, id)
		return s, err
	}
	s = &Statistic{
		StatisticID: id,
		Name:        name,
		Count:       int64(count),
	}
	cols, params := prepareColumnsWithValues(statisticColums)
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO statistics(%s) VALUES (%s)", cols, params), s.values()...); err != nil {
		return nil, err
	}
	return s, nil
}

func findStatistic(ctx context.Context, tx *sql.Tx, id string) (*Statistic, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}

	rows, err := tx.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM Statistics WHERE statistic_id=$1", strings.Join(statisticColums, ",")), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, nil
	}
	s, err := statisticFromRows(rows)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func statisticFromRows(rows *sql.Rows) (*Statistic, error) {
	var s Statistic
	err := rows.Scan(&s.StatisticID, &s.Name, &s.Count, &s.CreatedAt, &s.UpdatedAt)
	return &s, err
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
