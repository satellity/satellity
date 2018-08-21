package models

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"strings"
	"time"

	"github.com/go-pg/pg"
	"github.com/godiscourse/godiscourse/session"
	"github.com/godiscourse/godiscourse/uuid"
	"golang.org/x/crypto/bcrypt"
)

const users_DDL = `
CREATE TABLE IF NOT EXISTS users (
	user_id               VARCHAR(36) PRIMARY KEY,
	email                 VARCHAR(512) NOT NULL,
	username              VARCHAR(64) NOT NULL CHECK (username ~* '^[a-z0-9][a-z0-9_]{3,63}$'),
	nickname              VARCHAR(64) NOT NULL DEFAULT '',
	encrypted_password    VARCHAR(1024) NOT NULL,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX ON users ((LOWER(email)));
CREATE UNIQUE INDEX ON users ((LOWER(username)));
CREATE INDEX ON users (created_at);
`

type User struct {
	UserId            string    `sql:"user_id,pk"`
	Email             string    `sql:"email"`
	Username          string    `sql:"username"`
	Nickname          string    `sql:"nickname"`
	EncryptedPassword string    `sql:"encrypted_password"`
	CreatedAt         time.Time `sql:"created_at"`
	UpdatedAt         time.Time `sql:"updated_at"`

	SessionId string `sql:"-"`
}

var userCols = []string{"user_id", "email", "username", "nickname", "encrypted_password", "created_at", "updated_at"}

func CreateUser(ctx context.Context, email, username, nickname, password string, sessionSecret string) (*User, error) {
	t := time.Now()
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

	if err := validateEmailFormat(ctx, email); err != nil {
		return nil, err
	}
	if nickname == "" {
		nickname = username
	}
	password, err = validateAndEncryptPassword(ctx, password)
	if err != nil {
		return nil, err
	}

	user := &User{
		UserId:            uuid.NewV4().String(),
		Email:             email,
		Username:          username,
		Nickname:          nickname,
		EncryptedPassword: password,
		CreatedAt:         t,
		UpdatedAt:         t,
	}
	tx, err := session.Database(ctx).Begin()
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	defer tx.Rollback()

	if err := tx.Insert(user); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	sess, err := user.addSession(ctx, tx, sessionSecret)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	user.SessionId = sess.SessionId
	return user, nil
}

func FindUser(ctx context.Context, id string) (*User, error) {
	return findUserById(ctx, id)
}

func FindUserByUsernameOrEmail(ctx context.Context, q string) (*User, error) {
	user := &User{}
	q = strings.ToLower(strings.TrimSpace(q))
	if len(q) < 3 {
		return nil, nil
	}
	if err := session.Database(ctx).Model(user).Column(userCols...).Where("username = ? OR email = ?", q, q).Select(); err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return user, nil
}

func findUserById(ctx context.Context, id string) (*User, error) {
	user := &User{UserId: id}
	if err := session.Database(ctx).Model(user).Column(userCols...).WherePK().Select(); err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return user, nil
}

func validateAndEncryptPassword(ctx context.Context, password string) (string, error) {
	password = strings.TrimSpace(password)
	if len(password) < 8 {
		return password, session.PasswordTooSimpleError(ctx)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return password, session.ServerError(ctx, err)
	}
	return string(hashedPassword), nil
}
