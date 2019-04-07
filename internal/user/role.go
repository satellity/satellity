package user

import (
	"godiscourse/internal/configs"
)

// Role of an user, contains admin and member for now.
func (u *Data) Role() string {
	if configs.Operators[u.Email.String] {
		return "admin"
	}
	return "member"
}

// Name is nickname or username
func (u *Data) Name() string {
	if u.Nickname != "" {
		return u.Nickname
	}
	return u.Username
}

func (u *Data) isAdmin() bool {
	return u.Role() == "admin"
}
