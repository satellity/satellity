package user

import (
	"context"
	"godiscourse/internal/session"
	"net"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func validateEmailFormat(ctx context.Context, email string) error {
	if !emailRegexp.MatchString(email) {
		return session.InvalidEmailFormatError(ctx, email)
	}
	i := strings.LastIndexByte(email, '@')
	if _, err := net.LookupMX(email[i+1:]); err != nil {
		return session.InvalidEmailFormatError(ctx, email)
	}
	return nil
}

func validateAndEncryptPassword(ctx context.Context, password string) (string, error) {
	if len(password) < 8 {
		return password, session.PasswordTooSimpleError(ctx)
	}
	if len(password) > 64 {
		return password, session.BadDataError(ctx)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return password, session.ServerError(ctx, err)
	}
	return string(hashedPassword), nil
}
