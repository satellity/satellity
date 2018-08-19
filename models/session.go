package models

import (
	"context"
	"time"

	"github.com/go-pg/pg"
	"github.com/godiscourse/godiscourse/session"
	"github.com/godiscourse/godiscourse/uuid"
)

const sessions_DDL = `
CREATE TABLE IF NOT EXISTS sessions (
	session_id            VARCHAR(36) PRIMARY KEY,
	user_id               VARCHAR(36) NOT NULL,
	secret                VARCHAR(1024) NOT NULL,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE INDEX ON sessions (user_id);
CREATE INDEX ON sessions (created_at);
`

type Session struct {
	SessionId string    `sql:"session_id,pk"`
	UserId    string    `sql:"user_id"`
	Secret    string    `sql:"secret"`
	CreatedAt time.Time `sql:"created_at"`
}

var sessionCols = []string{"session_id", "user_id", "secret", "created_at"}

func (user *User) addSession(ctx context.Context, tx *pg.Tx, secret string) (*Session, error) {
	sess := &Session{
		SessionId: uuid.NewV4().String(),
		UserId:    user.UserId,
		Secret:    secret,
		CreatedAt: time.Now(),
	}
	err := tx.Insert(sess)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return sess, nil
}
