package user

import (
	"godiscourse/internal/configs"
)

// Role of an user, contains admin and member for now.
func (u *Model) Role() string {
	if configs.Operators[u.Email.String] {
		return "admin"
	}
	return "member"
}

// Name is nickname or username
func (u *Model) Name() string {
	if u.Nickname != "" {
		return u.Nickname
	}
	return u.Username
}

func (u *Model) IsAdmin() bool {
	return u.Role() == "admin"
}
