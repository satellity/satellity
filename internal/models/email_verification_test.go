package models

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailVerification(t *testing.T) {
	assert := assert.New(t)
	mctx := setupTestContext()
	defer mctx.database.Close()
	defer teardownTestContext(mctx)

	ev, err := CreateEmailVerification(mctx, "im.yuqlee@gmail.com", "testrecaptcha")
	assert.Nil(err)
	assert.NotNil(ev)

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Nil(err)
	public, err := x509.MarshalPKIXPublicKey(priv.Public())
	assert.Nil(err)
	user, err := VerifyEmailVerification(mctx, ev.VerificationID, ev.Code, "jason", "nopassword", hex.EncodeToString(public))
	assert.Nil(err)
	assert.NotNil(user)
	user, err = ReadUser(mctx, user.UserID)
	assert.Nil(err)
	assert.NotNil(user)
}
