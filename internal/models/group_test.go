package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupCRUD(t *testing.T) {
	assert := assert.New(t)
	mctx := setupTestContext()
	defer mctx.database.Close()
	defer teardownTestContext(mctx)

	user := createTestUser(mctx, "im.yuqlee@gmail.com", "username", "password")
	assert.NotNil(user)

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
		})
	}
}
