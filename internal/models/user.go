package models

import (
	"context"
	"database/sql"
	"fmt"
	"godiscourse/internal/configs"
	"godiscourse/internal/durable"
	"strings"
	"time"
)

const usersDDL = `
CREATE TABLE IF NOT EXISTS users (
	user_id               VARCHAR(36) PRIMARY KEY,
	email                 VARCHAR(512),
	username              VARCHAR(64) NOT NULL CHECK (username ~* '^[a-z0-9][a-z0-9_]{3,63}$'),
	nickname              VARCHAR(64) NOT NULL DEFAULT '',
	biography             VARCHAR(2048) NOT NULL DEFAULT '',
	encrypted_password    VARCHAR(1024),
	github_id             VARCHAR(1024) UNIQUE,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX ON users ((LOWER(email)));
CREATE UNIQUE INDEX ON users ((LOWER(username)));
CREATE INDEX ON users (created_at);
`

// User contains info of a register user
type User struct {
	UserID            string
	Email             sql.NullString
	Username          string
	Nickname          string
	Biography         string
	EncryptedPassword sql.NullString
	GithubID          sql.NullString
	CreatedAt         time.Time
	UpdatedAt         time.Time

	SessionID string
	isNew     bool
}

var userColumns = []string{"user_id", "email", "username", "nickname", "biography", "encrypted_password", "github_id", "created_at", "updated_at"}

func readUsersByIds(ctx context.Context, tx *sql.Tx, ids []string) ([]*User, error) {
	rows, err := tx.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM users WHERE user_id IN ('%s') LIMIT 100", strings.Join(userColumns, ","), strings.Join(ids, "','")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user, err := userFromRows(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func readUserSet(ctx context.Context, tx *sql.Tx, ids []string) (map[string]*User, error) {
	users, err := readUsersByIds(ctx, tx, ids)
	if err != nil {
		return nil, err
	}
	set := make(map[string]*User, 0)
	for _, u := range users {
		set[u.UserID] = u
	}
	return set, nil
}

// Role of an user, contains admin and member for now.
func (u *User) Role() string {
	if configs.Operators[u.Email.String] {
		return "admin"
	}
	return "member"
}

// Name is nickname or username
func (u *User) Name() string {
	if u.Nickname != "" {
		return u.Nickname
	}
	return u.Username
}

func (u *User) isAdmin() bool {
	return u.Role() == "admin"
}

// todo: move to statistics
func usersCount(ctx context.Context, tx *sql.Tx) (int64, error) {
	var count int64
	err := tx.QueryRowContext(ctx, "SELECT count(*) FROM users").Scan(&count)
	return count, err
}

func userFromRows(row durable.Row) (*User, error) {
	var u User
	err := row.Scan(&u.UserID, &u.Email, &u.Username, &u.Nickname, &u.Biography, &u.EncryptedPassword, &u.GithubID, &u.CreatedAt, &u.UpdatedAt)
	return &u, err
}
