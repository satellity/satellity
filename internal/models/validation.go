package models

import (
	"context"
	"net"
	"regexp"
	"satellity/internal/session"
	"strings"
)

var (
	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func validateEmailFormat(ctx context.Context, email string) error {
	if !emailRegexp.MatchString(email) {
		return session.InvalidEmailFormatError(ctx, email)
	}
	return nil
	i := strings.LastIndexByte(email, '@')
	if _, err := net.LookupMX(email[i+1:]); err != nil {
		return session.InvalidEmailFormatError(ctx, email)
	}
	return nil
}
