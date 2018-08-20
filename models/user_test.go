package models

import (
	"testing"

	"github.com/godiscourse/godiscourse/session"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer session.Database(ctx).Close()
	defer teardownTestContext(ctx)

	user, err := CreateUser(ctx, "im.yuqlee@gmailabcefgh.com", "username", "nickname", "password", "secret")
	assert.NotNil(err)
	assert.Nil(user)
	user, err = CreateUser(ctx, "im.yuqlee@gmail.com", "username", "nickname", "pass", "secret")
	assert.NotNil(err)
	assert.Nil(user)
	user, err = CreateUser(ctx, "im.yuqlee@gmail.com", "username", "nickname", "    pass     ", "secret")
	assert.NotNil(err)
	assert.Nil(user)
	user, err = CreateUser(ctx, "im.yuqlee@gmail.com", "username", "nickname", "password", "secret")
	assert.Nil(err)
	assert.NotNil(user)
	assert.NotEqual("", user.SessionId)
	new, err := FindUser(ctx, user.UserId)
	assert.Nil(err)
	assert.NotNil(new)
	assert.Equal(user.Username, new.Username)
	assert.Equal(user.Nickname, new.Nickname)
	err = bcrypt.CompareHashAndPassword([]byte(new.EncryptedPassword), []byte("password"))
	assert.Nil(err)
	new, err = FindUserByUsernameOrEmail(ctx, "None")
	assert.Nil(err)
	assert.Nil(new)
	new, err = FindUserByUsernameOrEmail(ctx, "im.yuqlee@Gmail.com")
	assert.Nil(err)
	assert.NotNil(new)
	new, err = FindUserByUsernameOrEmail(ctx, "UserName")
	assert.Nil(err)
	assert.NotNil(new)
	new, err = FindUserByUsernameOrEmail(ctx, "im.yuqlee@Gmail.com")
	assert.Nil(err)
	assert.NotNil(new)
}
