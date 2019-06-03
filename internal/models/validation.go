package models

import (
	"context"
	"godiscourse/internal/session"
	"net"
	"regexp"
	"strings"
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
