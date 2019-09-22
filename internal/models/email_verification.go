package models

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"fmt"
	"satellity/internal/clouds"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

const emailVerificationDDL = `
CREATE TABLE IF NOT EXISTS email_verifications (
	verification_id        VARCHAR(36) PRIMARY KEY,
	email                  VARCHAR(512),
	code                   VARCHAR(512),
	created_at             TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS email_verifications_email_code_createdx ON email_verifications (email, code, created_at DESC);
CREATE INDEX IF NOT EXISTS email_verifications_createdx ON email_verifications (created_at DESC);
`

// EmailVerification verify email
type EmailVerification struct {
	VerificationID string
	Email          string
	Code           string
	CreatedAt      time.Time
}

var emailVerificationColumns = []string{"verification_id", "email", "code", "created_at"}

func (e *EmailVerification) values() []interface{} {
	return []interface{}{e.VerificationID, e.Email, e.Code, e.CreatedAt}
}

func emailVerificationFromRows(row durable.Row) (*EmailVerification, error) {
	var ev EmailVerification
	err := row.Scan(&ev.VerificationID, &ev.Email, &ev.Code, &ev.CreatedAt)
	return &ev, err
}

// CreateEmailVerification create an email verification
func CreateEmailVerification(mctx *Context, email, recaptcha string) (*EmailVerification, error) {
	ctx := mctx.context

	code, err := generateVerificationCode(ctx)
	if err != nil {
		return nil, err
	}

	success, err := clouds.VerifyRecaptcha(ctx, recaptcha)
	if err != nil {
		return nil, session.ServerError(ctx, err)
	} else if !success {
		return nil, session.RecaptchaVerifyError(ctx)
	}
	ev := &EmailVerification{
		VerificationID: uuid.Must(uuid.NewV4()).String(),
		Email:          strings.TrimSpace(email),
		Code:           code,
		CreatedAt:      time.Now(),
	}

	err = mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		_, err = tx.ExecContext(ctx, "DELETE FROM email_verifications WHERE created_at<$1", time.Now().Add(-24*time.Hour))
		if err != nil {
			return err
		}
		columns, params := durable.PrepareColumnsWithValues(emailVerificationColumns)
		query := fmt.Sprintf("INSERT INTO email_verifications(%s) VALUES (%s)", columns, params)
		_, err := tx.ExecContext(ctx, query, ev.values()...)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return ev, nil
}

// VerifyEmailVerification verify an email verification
func VerifyEmailVerification(mctx *Context, verificationID, code, username, password, sessionSecret string) (*User, error) {
	ctx := mctx.context

	var user *User
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		ev, err := findEmailVerification(ctx, tx, verificationID)
		if err != nil || ev == nil {
			return err
		}
		if ev.Code != code {
			ev, err = findEmailVerificationByEmailAndCode(ctx, tx, ev.Email, code)
			if err != nil || ev == nil {
				return err
			}
		}
		if ev.CreatedAt.Add(time.Hour * 24).Before(time.Now()) {
			return session.VerificationCodeInvalidError(ctx)
		}
		_, err = tx.ExecContext(ctx, "DELETE FROM email_verifications WHERE verification_id=$1", ev.VerificationID)
		if err != nil {
			return err
		}

		password, err = validateAndEncryptPassword(ctx, password)
		if err != nil {
			return err
		}
		t := time.Now()
		user = &User{
			UserID:            uuid.Must(uuid.NewV4()).String(),
			Email:             sql.NullString{String: ev.Email, Valid: true},
			Username:          username,
			Nickname:          username,
			EncryptedPassword: sql.NullString{String: password, Valid: true},
			CreatedAt:         t,
			UpdatedAt:         t,
		}
		columns, params := durable.PrepareColumnsWithValues(userColumns)
		_, err = tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO users(%s) VALUES (%s)", columns, params), user.values()...)
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
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return nil, session.TransactionError(ctx, err)
	} else if user == nil {
		return nil, session.VerificationCodeInvalidError(ctx)
	}
	return user, nil
}

func findEmailVerification(ctx context.Context, tx *sql.Tx, id string) (*EmailVerification, error) {
	row := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM email_verifications WHERE verification_id=$1", strings.Join(emailVerificationColumns, ",")), id)
	ev, err := emailVerificationFromRows(row)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return ev, nil
}

func findEmailVerificationByEmailAndCode(ctx context.Context, tx *sql.Tx, email, code string) (*EmailVerification, error) {
	row := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM email_verifications WHERE email=$1 AND code=$2 ORDER BY email,code,created_at DESC LIMIT 1", strings.Join(emailVerificationColumns, ",")), email, code)
	ev, err := emailVerificationFromRows(row)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return ev, nil
}

func generateVerificationCode(ctx context.Context) (string, error) {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return "", session.ServerError(ctx, err)
	}
	c := binary.LittleEndian.Uint64(b[:]) % 10000
	if c < 1000 {
		c = 1000 + c
	}
	return fmt.Sprint(c), nil
}
