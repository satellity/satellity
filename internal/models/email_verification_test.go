package models

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
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

	pub, _, err := ed25519.GenerateKey(rand.Reader)
	assert.Nil(err)
	user, err := VerifyEmailVerification(ctx, ev.VerificationID, ev.Code, "jason", "nopassword", base64.RawURLEncoding.EncodeToString(pub))
	assert.Nil(err)
	assert.NotNil(user)
	user, err = ReadUser(ctx, user.UserID)
	assert.Nil(err)
	assert.NotNil(user)
}
