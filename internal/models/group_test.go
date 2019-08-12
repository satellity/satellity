package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGroupCRUD(t *testing.T) {
	assert := assert.New(t)
	mctx := setupTestContext()
	defer mctx.database.Close()
	defer teardownTestContext(mctx)

	user := createTestUser(mctx, "im.yuqlee@gmail.com", "username", "password")
	assert.NotNil(user)
	jason := createTestUser(mctx, "validfake@gmail.com", "usernamex", "passwordx")
	assert.NotNil(jason)

	groupCases := []struct {
		name        string
		description string
		valid       bool
	}{
		{"iv", "invalid group name", false},
		{"group", "valid group name", true},
	}

	for _, tc := range groupCases {
		t.Run(fmt.Sprintf("group name %s", tc.name), func(t *testing.T) {
			if !tc.valid {
				group, err := user.CreateGroup(mctx, tc.name, tc.description)
				assert.NotNil(err)
				assert.Nil(group)
				return
			}

			group, err := user.CreateGroup(mctx, tc.name, tc.description)
			assert.Nil(err)
			assert.NotNil(group)

			new, err := ReadGroup(mctx, uuid.Must(uuid.NewV4()).String(), nil)
			assert.Nil(err)
			assert.Nil(new)
			new, err = ReadGroup(mctx, group.GroupID, nil)
			assert.Nil(err)
			assert.NotNil(new)
			users, err := new.Participants(mctx, nil, time.Now(), "100")
			assert.Nil(err)
			assert.Len(users, 1)
			groups, err := user.ReadGroups(mctx)
			assert.Nil(err)
			assert.Len(groups, 1)
			groups, err = ReadGroups(mctx)
			assert.Nil(err)
			assert.Len(groups, 1)
			groups, err = user.RelatedGroups(mctx, 100)
			assert.Nil(err)
			assert.Len(groups, 1)

			name := "new" + tc.name
			description := "new" + tc.description
			group, err = user.UpdateGroup(mctx, group.GroupID, name, description)
			assert.Nil(err)
			assert.NotNil(group)
			assert.Equal(name, group.Name)
			assert.Equal(description, group.Description)
			new, err = ReadGroup(mctx, group.GroupID, nil)
			assert.Nil(err)
			assert.NotNil(new)
			assert.Equal(name, new.Name)
			assert.Equal(description, new.Description)

			err = jason.JoinGroup(mctx, group.GroupID, "INVALID")
			assert.NotNil(err)
			err = jason.JoinGroup(mctx, group.GroupID, ParticipantRoleMember)
			assert.Nil(err)
			new, _ = ReadGroup(mctx, group.GroupID, nil)
			assert.Equal(int64(2), new.UsersCount)
			users, err = group.Participants(mctx, nil, time.Now(), "100")
			assert.Len(users, 2)
			err = jason.ExitGroup(mctx, group.GroupID)
			assert.Nil(err)
			new, _ = ReadGroup(mctx, group.GroupID, nil)
			assert.Equal(int64(1), new.UsersCount)
			users, err = group.Participants(mctx, nil, time.Now(), "100")
			assert.Len(users, 1)
		})
	}
}
