package models

import (
	"testing"

	"github.com/godiscourse/godiscourse/session"
	"github.com/stretchr/testify/assert"
)

func TestTopicCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer session.Database(ctx).Close()
	defer teardownTestContext(ctx)

	user := createTestUser(ctx, "im.yuqlee@gmail.com", "username", "password")
	assert.NotNil(user)
	category, _ := CreateCategory(ctx, "name", "Description")
	assert.NotNil(category)
	topic, err := user.CreateTopic(ctx, "title", "body", category.CategoryID)
	assert.Nil(err)
	assert.NotNil(topic)
	category, _ = ReadCategory(ctx, category.CategoryID)
	assert.NotNil(category)
	assert.Equal(topic.TopicId, category.LastTopicID.String)
	assert.Equal(1, category.TopicsCount)
}
