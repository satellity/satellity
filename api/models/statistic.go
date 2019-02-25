package models

import (
	"context"
	"crypto/md5"
	"io"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/godiscourse/godiscourse/api/session"
	uuid "github.com/satori/go.uuid"
)

const STATISTIC_ID = "540cbd3c-f4eb-479c-bcd8-b5629af57267"
const statisticsDDL = `
CREATE TABLE IF NOT EXISTS statistics (
	statistic_id          VARCHAR(36) PRIMARY KEY,
	name                  VARCHAR(512) NOT NULL,
	count                 BIGINT NOT NULL DEFAULT 0,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
`

type Statistic struct {
	StatisticID string    `sql:"statistic_id,pk"`
	Name        string    `sql:"name,notnull"`
	Count       int64     `sql:"count,notnull"`
	CreatedAt   time.Time `sql:"created_at"`
	UpdatedAt   time.Time `sql:"updated_at"`
}

var statisticColums = []string{"statistic_id", "name", "count", "created_at", "updated_at"}

func upsertStatistic(ctx context.Context, name string) (*Statistic, error) {
	id, _ := generateStatisticId(STATISTIC_ID, name)
	s, err := readStatistic(ctx, id)
	if err != nil {
		return nil, err
	}
	err = session.Database(ctx).RunInTransaction(func(tx *pg.Tx) error {
		var count int
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
		if s != nil {
			s.Count = int64(count)
			return tx.Update(s)
		}
		s = &Statistic{
			StatisticID: id,
			Name:        name,
			Count:       int64(count),
		}
		return tx.Insert(s)
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return s, nil
}

func readStatistic(ctx context.Context, id string) (*Statistic, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}
	s := &Statistic{StatisticID: id}
	if err := session.Database(ctx).Model(s).WherePK().Select(); err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return s, nil
}

// BeforeInsert hook insert
func (s *Statistic) BeforeInsert(db orm.DB) error {
	s.CreatedAt = time.Now()
	s.UpdatedAt = s.CreatedAt
	return nil
}

// BeforeUpdate hook update
func (s *Statistic) BeforeUpdate(db orm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}

func generateStatisticId(Id, name string) (string, error) {
	h := md5.New()
	io.WriteString(h, Id)
	io.WriteString(h, name)
	sum := h.Sum(nil)
	sum[6] = (sum[6] & 0x0f) | 0x30
	sum[8] = (sum[8] & 0x3f) | 0x80
	id, err := uuid.FromBytes(sum)
	return id.String(), err
}
