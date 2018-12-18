package models

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/godiscourse/godiscourse/api/session"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer session.Database(ctx).Close()
	defer teardownTestContext(ctx)

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Nil(err)
	public, err := x509.MarshalPKIXPublicKey(priv.Public())
	assert.Nil(err)
	user, err := CreateUser(ctx, "im.yuqlee@gmailabcefgh.com", "username", "nickname", "", "password", hex.EncodeToString(public))
	assert.NotNil(err)
	assert.Nil(user)
	user, err = CreateUser(ctx, "im.yuqlee@gmail.com", "username", "nickname", "", "pass", hex.EncodeToString(public))
	assert.NotNil(err)
	assert.Nil(user)
	user, err = CreateUser(ctx, "im.yuqlee@gmail.com", "username", "nickname", "", "    pass     ", hex.EncodeToString(public))
	assert.NotNil(err)
	assert.Nil(user)
	user, err = CreateUser(ctx, "im.yuqlee@gmail.com", "username", "nickname", "", "password", hex.EncodeToString(public))
	assert.Nil(err)
	assert.NotNil(user)
	assert.NotEqual("", user.SessionID)
	new, err := ReadUser(ctx, user.UserID)
	assert.Nil(err)
	assert.NotNil(new)
	assert.Equal(user.Username, new.Username)
	assert.Equal(user.Nickname, new.Nickname)
	err = bcrypt.CompareHashAndPassword([]byte(new.EncryptedPassword.String), []byte("password"))
	assert.Nil(err)
	new, err = ReadUserByUsernameOrEmail(ctx, "None")
	assert.Nil(err)
	assert.Nil(new)
	new, err = ReadUserByUsernameOrEmail(ctx, "im.yuqlee@Gmail.com")
	assert.Nil(err)
	assert.NotNil(new)
	new, err = ReadUserByUsernameOrEmail(ctx, "UserName")
	assert.Nil(err)
	assert.NotNil(new)
	new, err = ReadUserByUsernameOrEmail(ctx, "im.yuqlee@Gmail.com")
	assert.Nil(err)
	assert.NotNil(new)
	new, err = CreateSession(ctx, "im.yuqlee@Gmail.com", "password", hex.EncodeToString(public))
	assert.Nil(err)
	assert.NotNil(new)
	assert.Equal("username", user.Username)
	assert.Equal("member", user.Role())

	sess, err := readSession(ctx, new.UserID, new.SessionID)
	assert.Nil(err)
	assert.NotNil(sess)
	sess, err = readSession(ctx, uuid.Must(uuid.NewV4()).String(), new.SessionID)
	assert.Nil(err)
	assert.Nil(sess)

	claims := &jwt.MapClaims{
		"uid": new.UserID,
		"sid": new.SessionID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	ss, err := token.SignedString(priv)
	assert.Nil(err)
	new, err = AuthenticateUser(ctx, ss)
	assert.Nil(err)
	assert.NotNil(new)
	err = new.UpdateProfile(ctx, "Jason", "")
	assert.Nil(err)
	assert.Equal("Jason", new.Name())
	new, err = ReadUserByUsernameOrEmail(ctx, "UserName")
	assert.Nil(err)
	assert.NotNil(new)
	assert.Equal("Jason", new.Name())
	users, err := ReadUsers(ctx, time.Time{})
	assert.Nil(err)
	assert.Len(users, 1)
}

func createTestUser(ctx context.Context, email, username, password string) *User {
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	public, _ := x509.MarshalPKIXPublicKey(priv.Public())
	user, _ := CreateUser(ctx, email, username, "nickname", "", password, hex.EncodeToString(public))
	return user
}
