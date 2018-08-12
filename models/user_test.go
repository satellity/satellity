package models

import (
	"testing"

	"github.com/godiscourse/godiscourse/session"
	"github.com/stretchr/testify/assert"
)

func TestUserCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer session.Database(ctx).Close()
	defer teardownTestContext(ctx)

	user, err := CreateUser(ctx, "username", "nickname")
	assert.Nil(err)
	assert.NotNil(user)
	new, err := FindUser(ctx, user.UserId)
	assert.Nil(err)
	assert.NotNil(new)
	assert.Equal(user.Username, new.Username)
	assert.Equal(user.Nickname, new.Nickname)
}
