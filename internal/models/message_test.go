package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMessageCRUD(t *testing.T) {
	assert := assert.New(t)
	mctx := setupTestContext()
	defer mctx.database.Close()
	defer teardownTestContext(mctx)

	user := createTestUser(mctx, "im.yuqlee@gmail.com", "username", "password")
	assert.NotNil(user)
	jade := createTestUser(mctx, "im.jadeydi@gmail.com", "username1", "password")
	assert.NotNil(jade)
	group, err := user.CreateGroup(mctx, "Group Name", "Group Description")
	assert.Nil(err)
	assert.NotNil(group)
	users, err := group.Participants(mctx, nil, time.Now(), "100")
	assert.Len(users, 1)

	messageCases := []struct {
		Body  string
		Valid bool
	}{
		{"Message Body", true},
	}

	for _, tc := range messageCases {
		t.Run(fmt.Sprintf("Message %s", tc.Body), func(t *testing.T) {
			message, err := user.CreateMessage(mctx, uuid.Must(uuid.NewV4()).String(), tc.Body, "")
			assert.NotNil(err)
			assert.Nil(message)
			message, err = user.CreateMessage(mctx, group.GroupID, tc.Body, "")
			assert.Nil(err)
			assert.NotNil(message)
			new, err := ReadMessage(mctx, message.MessageID)
			assert.Nil(err)
			assert.NotNil(new)
			new, err = ReadMessage(mctx, uuid.Must(uuid.NewV4()).String())
			assert.Nil(err)
			assert.Nil(new)
			messages, err := group.ReadMessages(mctx, time.Now())
			assert.Nil(err)
			assert.Len(messages, 1)
			new, err = jade.UpdateMessage(mctx, message.MessageID, tc.Body)
			assert.NotNil(err)
			assert.Nil(new)
			new, err = user.UpdateMessage(mctx, message.MessageID, "New"+tc.Body)
			assert.Nil(err)
			assert.NotNil(new)
			sub, err := user.CreateMessage(mctx, group.GroupID, tc.Body, new.MessageID)
			assert.Nil(err)
			assert.NotNil(sub)
			messages, err = group.ReadMessages(mctx, time.Now())
			assert.Nil(err)
			assert.Len(messages, 2)
			err = jade.DeleteMessage(mctx, message.MessageID)
			assert.NotNil(err)
			err = user.DeleteMessage(mctx, message.MessageID)
			assert.Nil(err)
			new, err = ReadMessage(mctx, message.MessageID)
			assert.Nil(err)
			assert.Nil(new)
		})
	}
}
