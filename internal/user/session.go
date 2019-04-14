package user

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"database/sql"
	"encoding/hex"
	"fmt"
	"godiscourse/internal/durable"
	"godiscourse/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

const sessionsDDL = `
CREATE TABLE IF NOT EXISTS sessions (
	session_id            VARCHAR(36) PRIMARY KEY,
	user_id               VARCHAR(36) NOT NULL,
	secret                VARCHAR(1024) NOT NULL,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE INDEX ON sessions (user_id);
`

// Session contains user's current login infomation
type Session struct {
	SessionID string    `sql:"session_id,pk"`
	UserID    string    `sql:"user_id"`
	Secret    string    `sql:"secret"`
	CreatedAt time.Time `sql:"created_at"`
}

var SessionColumns = []string{"session_id", "user_id", "secret", "created_at"}

func (s *Session) sessionValues() []interface{} {
	return []interface{}{s.SessionID, s.UserID, s.Secret, s.CreatedAt}
}

func checkSecret(ctx context.Context, sessionSecret string) error {
	data, err := hex.DecodeString(sessionSecret)
	if err != nil {
		return session.BadDataError(ctx)
	}
	public, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return session.BadDataError(ctx)
	}
	switch public.(type) {
	case *ecdsa.PublicKey:
	default:
		return session.BadDataError(ctx)
	}
	return nil
}

func (d *Model) addSession(ctx context.Context, tx *sql.Tx, secret string) (*Session, error) {
	s := &Session{
		SessionID: uuid.Must(uuid.NewV4()).String(),
		UserID:    d.UserID,
		Secret:    secret,
		CreatedAt: time.Now(),
	}

	cols, params := durable.PrepareColumnsWithValues(SessionColumns)
	_, err := tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO sessions(%s) VALUES(%s)", cols, params), s.sessionValues()...)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return s, nil
}

func readSession(ctx context.Context, tx *sql.Tx, uid, sid string) (*Session, error) {
	if id, _ := uuid.FromString(uid); id.String() == uuid.Nil.String() {
		return nil, nil
	}
	if id, _ := uuid.FromString(sid); id.String() == uuid.Nil.String() {
		return nil, nil
	}

	row := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM sessions WHERE user_id=$1 AND session_id=$2", strings.Join(SessionColumns, ",")), uid, sid)
	s, err := sessionFromRows(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

func sessionFromRows(row durable.Row) (*Session, error) {
	var s Session
	err := row.Scan(&s.SessionID, &s.UserID, &s.Secret, &s.CreatedAt)
	return &s, err
}
