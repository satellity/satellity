package models

import (
	"context"
	"time"

	"github.com/go-pg/pg"
	"github.com/godiscourse/godiscourse/session"
	"github.com/godiscourse/godiscourse/uuid"
)

const users_DDL = `
CREATE TABLE IF NOT EXISTS users (
	user_id			VARCHAR(36) PRIMARY KEY,
	username		VARCHAR(64) NOT NULL CHECK (username ~* '^[a-z0-9][a-z0-9_]{3,63}$'),
	nickname		VARCHAR(64) NOT NULL DEFAULT '',
	created_at	TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at	TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX ON users ((LOWER(username)));
CREATE INDEX ON users (created_at);
`

type User struct {
	UserId    string    `sql:"user_id,pk"`
	Username  string    `sql:"username"`
	Nickname  string    `sql:"nickname"`
	CreatedAt time.Time `sql:"created_at"`
	UpdatedAt time.Time `sql:"updated_at"`
}

var userCols = []string{"user_id", "username", "nickname", "created_at", "updated_at"}

func CreateUser(ctx context.Context, username, nickname string) (*User, error) {
	t := time.Now()

	if nickname == "" {
		nickname = username
	}
	user := &User{
		UserId:    uuid.NewV4().String(),
		Username:  username,
		Nickname:  nickname,
		CreatedAt: t,
		UpdatedAt: t,
	}
	if err := session.Database(ctx).Insert(user); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return user, nil
}

func FindUser(ctx context.Context, id string) (*User, error) {
	return findUserById(ctx, id)
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
