package models

import "time"

const users_DDL = `
`

type User struct {
	UserId    string
	Username  string
	Nickname  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
