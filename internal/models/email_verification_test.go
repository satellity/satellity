package models

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailVerification(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer teardownTestContext(ctx)

	ev, err := CreateEmailVerification(ctx, "USER", "im.yuqlee@gmail.com", "testrecaptcha")
	assert.Nil(err)
	assert.NotNil(ev)

	public, _, err := ed25519.GenerateKey(rand.Reader)
	assert.Nil(err)
	user, err := VerifyEmailVerification(ctx, ev.VerificationID, ev.Code, "nopassword", hex.EncodeToString(public))
	assert.Nil(err)
	assert.NotNil(user)
	user, err = ReadUser(ctx, user.UserID)
	assert.Nil(err)
	assert.NotNil(user)
}
