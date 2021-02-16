package models

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
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

	err = session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "DELETE FROM sessions WHERE session_id IN (SELECT session_id FROM sessions WHERE user_id=$1 ORDER BY created_at DESC OFFSET 5)", user.UserID)
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

func (user *User) addSession(ctx context.Context, tx pgx.Tx, secret string) (*Session, error) {
	s := &Session{
		SessionID: uuid.Must(uuid.NewV4()).String(),
		UserID:    user.UserID,
		Secret:    secret,
		CreatedAt: time.Now(),
	}

	rows := [][]interface{}{s.values()}
	_, err := tx.CopyFrom(ctx, pgx.Identifier{"sessions"}, sessionColumns, pgx.CopyFromRows(rows))
	return s, err
}

func readSession(ctx context.Context, tx pgx.Tx, uid, sid string) (*Session, error) {
	if id, _ := uuid.FromString(uid); id.String() == uuid.Nil.String() {
		return nil, nil
	}
	if id, _ := uuid.FromString(sid); id.String() == uuid.Nil.String() {
		return nil, nil
	}

	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM sessions WHERE user_id=$1 AND session_id=$2", strings.Join(sessionColumns, ",")), uid, sid)
	s, err := sessionFromRows(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return s, err
}
