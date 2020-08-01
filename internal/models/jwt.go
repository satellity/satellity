package models

import (
	"crypto/ed25519"

	jwt "github.com/dgrijalva/jwt-go"
)

var Ed25519SigningMethod *EdDSASigningMethod

func init() {
	Ed25519SigningMethod = &EdDSASigningMethod{}
	jwt.RegisterSigningMethod("EdDSA", func() jwt.SigningMethod {
		return Ed25519SigningMethod
	})
}

type EdDSASigningMethod struct{}

func (sm *EdDSASigningMethod) Verify(signingString, signature string, key interface{}) error {
	switch k := key.(type) {
	case ed25519.PublicKey:
		sig, err := jwt.DecodeSegment(signature)
		if err != nil {
			return err
		}
		if !ed25519.Verify(k, []byte(signingString), sig) {
			return jwt.ErrECDSAVerification
		}
	default:
		return jwt.ErrInvalidKeyType
	}
	return nil
}

func (sm *EdDSASigningMethod) Sign(signingString string, key interface{}) (string, error) {
	switch k := key.(type) {
	case ed25519.PrivateKey:
		sig := ed25519.Sign(k, []byte(signingString))
		return jwt.EncodeSegment(sig), nil
	default:
		return "", jwt.ErrInvalidKeyType
	}
}

func (sm *EdDSASigningMethod) Alg() string {
	return "EdDSA"
}
