package model

import (
	"database/sql"
	"time"
)

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

type UserInfo struct {
	Email         string
	Username      string
	Nickname      string
	Biography     string
	Password      string
	SessionSecret string
}