package models

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"database/sql"
	"encoding/hex"
	"fmt"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Session contains user's current login infomation
type Session struct {
	SessionID string    `sql:"session_id,pk"`
	UserID    string    `sql:"user_id"`
	Secret    string    `sql:"secret"`
	CreatedAt time.Time `sql:"created_at"`
}

var sessionColumns = []string{"session_id", "user_id", "secret", "created_at"}

func (s *Session) values() []interface{} {
	return []interface{}{s.SessionID, s.UserID, s.Secret, s.CreatedAt}
}

func sessionFromRows(row durable.Row) (*Session, error) {
	var s Session
	err := row.Scan(&s.SessionID, &s.UserID, &s.Secret, &s.CreatedAt)
	return &s, err
}

// CreateSession create a new user session
func CreateSession(mctx *Context, identity, password, sessionSecret string) (*User, error) {
	ctx := mctx.context
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

	user, err := ReadUserByUsernameOrEmail(mctx, identity)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, session.IdentityNonExistError(ctx)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword.String), []byte(password)); err != nil {
		return nil, session.InvalidPasswordError(ctx)
	}

	err = mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM sessions WHERE session_id IN (SELECT session_id FROM sessions WHERE user_id=$1 ORDER BY created_at DESC OFFSET 5)", user.UserID)
		if err != nil {
			return err
		}
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

	cols, posits := durable.PrepareColumnsWithParams(sessionColumns)
	stmt, err := tx.PrepareContext(ctx, fmt.Sprintf("INSERT INTO sessions(%s) VALUES(%s)", cols, posits))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, s.values()...)
	return s, nil
}

func readSession(ctx context.Context, tx *sql.Tx, uid, sid string) (*Session, error) {
	if id, _ := uuid.FromString(uid); id.String() == uuid.Nil.String() {
		return nil, nil
	}
	if id, _ := uuid.FromString(sid); id.String() == uuid.Nil.String() {
		return nil, nil
	}

	row := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM sessions WHERE user_id=$1 AND session_id=$2", strings.Join(sessionColumns, ",")), uid, sid)
	s, err := sessionFromRows(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

const sessionsDDL = `
CREATE TABLE IF NOT EXISTS sessions (
	session_id            VARCHAR(36) PRIMARY KEY,
	user_id               VARCHAR(36) NOT NULL,
	secret                VARCHAR(1024) NOT NULL,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE INDEX ON sessions (user_id);
`

const dropSessionsDDL = `DROP TABLE IF EXISTS sessions;`
