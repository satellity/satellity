package models

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/godiscourse/godiscourse/api/session"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
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

var sessionCols = []string{"session_id", "user_id", "secret", "created_at"}

func (s *Session) values() []interface{} {
	return []interface{}{s.SessionID, s.UserID, s.Secret, s.CreatedAt}
}

// CreateSession create a new user session
func CreateSession(ctx context.Context, identity, password, sessionSecret string) (*User, error) {
	data, err := hex.DecodeString(sessionSecret)
	if err != nil {
		return nil, session.BadDataError(ctx)
	}
	public, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, session.BadDataError(ctx)
	}
	switch public.(type) {
	case *ecdsa.PublicKey:
	default:
		return nil, session.BadDataError(ctx)
	}

	user, err := ReadUserByUsernameOrEmail(ctx, identity)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, session.IdentityNonExistError(ctx)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword.String), []byte(password)); err != nil {
		return nil, session.InvalidPasswordError(ctx)
	}

	err = runInTransaction(ctx, func(tx *sql.Tx) error {
		s, err := user.addSession(ctx, tx, sessionSecret)
		if err != nil {
			return err
		}
		user.SessionID = s.SessionID
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return user, nil
}

func (user *User) addSession(ctx context.Context, tx *sql.Tx, secret string) (*Session, error) {
	s := &Session{
		SessionID: uuid.Must(uuid.NewV4()).String(),
		UserID:    user.UserID,
		Secret:    secret,
		CreatedAt: time.Now(),
	}

	cols, params := prepareColumnsWithValues(sessionCols)
	_, err := tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO sessions(%s) VALUES(%s)", cols, params), s.values()...)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return s, nil
}

func readSession(ctx context.Context, tx *sql.Tx, uid, sid string) (*Session, error) {
	if id, _ := uuid.FromString(uid); id.String() == uuid.Nil.String() {
		return nil, nil
	}
	if id, _ := uuid.FromString(uid); id.String() == uuid.Nil.String() {
		return nil, nil
	}

	rows, err := tx.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM sessions WHERE user_id=$1 AND session_id=$2", strings.Join(sessionCols, ",")), uid, sid)
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
	s, err := sessionFromRows(rows)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func sessionFromRows(rows *sql.Rows) (*Session, error) {
	var s Session
	err := rows.Scan(&s.SessionID, &s.UserID, &s.Secret, &s.CreatedAt)
	return &s, err
}
