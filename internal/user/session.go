package user

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"godiscourse/internal/session"
)

func checkSecret(ctx context.Context, sessionSecret string) error {
	data, err := hex.DecodeString(sessionSecret)
	if err != nil {
		return session.BadDataError(ctx)
	}
	public, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return session.BadDataError(ctx)
	}
	switch public.(type) {
	case *ecdsa.PublicKey:
	default:
		return session.BadDataError(ctx)
	}
	return nil
}
