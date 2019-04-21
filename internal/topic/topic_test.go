package topic

import (
	"context"
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"

	"godiscourse/internal/category"
	"godiscourse/internal/user"
)

func TestTopicCRUD(t *testing.T) {
	t.Skip()
	assert := assert.New(t)
	ctx := context.Background()

	userMock := user.NewMock()
	// todo: categoryMock

	user, err := userMock.Create(ctx, &user.Params{
		Email:    "im.yuqlee@gmail.com",
		Username: "username",
		Password: "password",
	})
	assert.Nil(err)
	assert.NotNil(user)
	category, _ := category.Create(ctx, &category.Params{
		Name:        "name",
		Alias:       "alias",
		Description: "Description",
		Position:    0,
	})
	assert.NotNil(category)

	topicCases := []struct {
		title         string
		body          string
		categoryID    string
		topicsCount   int
		commentsCount int
		valid         bool
	}{
		{"", "body", category.CategoryID, 0, 0, false},
		{"title", "body", uuid.Must(uuid.NewV4()).String(), 0, 0, false},
		{"title", "body", category.CategoryID, 1, 0, true},
		{"title2", "body", category.CategoryID, 2, 0, true},
	}

	for _, tc := range topicCases {
		t.Run(fmt.Sprintf("topic title %s", tc.title), func(t *testing.T) {
			// todo: integration tests
		})
	}

	user, err = userMock.Create(ctx, &user.Params{
		Email:    "im.jadeydi@gmail.com",
		Username: "usernamex",
		Password: "password",
	})
	assert.Nil(err)
	assert.NotNil(user)
	category, _ := category.Create(ctx, &category.Params{
		Name:        "new name",
		Alias:       "new alias",
		Description: "New Description",
		Position:    2,
	})
	assert.NotNil(category)
	topicCases = []struct {
		title         string
		body          string
		categoryID    string
		topicsCount   int
		commentsCount int
		valid         bool
	}{
		{"title", "body", category.CategoryID, 1, 0, true},
		{"title2", "body", category.CategoryID, 2, 0, true},
	}

	for _, tc := range topicCases {
		t.Run(fmt.Sprintf("topic title %s", tc.title), func(t *testing.T) {
			// todo: integration tests
		})
	}
}
