package models

import (
	"testing"
	"time"

	"github.com/godiscourse/godiscourse/api/session"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestTopicCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer session.Database(ctx).Close()
	defer teardownTestContext(ctx)

	user := createTestUser(ctx, "im.yuqlee@gmail.com", "username", "password")
	assert.NotNil(user)
	category, _ := CreateCategory(ctx, "name", "Description", 0)
	assert.NotNil(category)
	topic, err := user.CreateTopic(ctx, "title", "body", category.CategoryID)
	assert.Nil(err)
	assert.NotNil(topic)
	category, _ = ReadCategory(ctx, category.CategoryID)
	assert.NotNil(category)
	assert.Equal(topic.TopicID, category.LastTopicID.String)
	assert.Equal(1, category.TopicsCount)
	topic, err = ReadTopic(ctx, topic.TopicID)
	assert.Nil(err)
	assert.NotNil(topic)
	topics, err := ReadTopics(ctx, time.Time{})
	assert.Nil(err)
	assert.Len(topics, 1)
	topics, err = user.ReadTopics(ctx, time.Time{})
	assert.Nil(err)
	assert.Len(topics, 1)
	topics, err = category.ReadTopics(ctx, time.Time{})
	assert.Nil(err)
	assert.Len(topics, 1)

	user = createTestUser(ctx, "im.jadeydi@gmail.com", "usernamex", "password")
	assert.NotNil(user)
	topic, err = user.CreateTopic(ctx, "title", "body", category.CategoryID)
	assert.Nil(err)
	assert.NotNil(topic)
	topics, err = ReadTopics(ctx, time.Time{})
	assert.Nil(err)
	assert.Len(topics, 2)
	topics, err = user.ReadTopics(ctx, time.Time{})
	assert.Nil(err)
	assert.Len(topics, 1)
	topic, err = user.UpdateTopic(ctx, topic.TopicID, "hell", "orld")
	assert.Nil(err)
	assert.NotNil(topic)
	assert.Equal("hell", topic.Title)
	assert.Equal("orld", topic.Body)
	topic, err = user.UpdateTopic(ctx, topic.TopicID, "", "orld orld")
	assert.Nil(err)
	assert.NotNil(topic)
	assert.Equal("hell", topic.Title)
	assert.Equal("orld orld", topic.Body)
	new, err := user.UpdateTopic(ctx, uuid.Must(uuid.NewV4()).String(), "hell", "orld")
	assert.NotNil(err)
	assert.Nil(new)
	u := &User{UserID: uuid.Must(uuid.NewV4()).String()}
	new, err = u.UpdateTopic(ctx, topic.TopicID, "hell", "orld")
	assert.NotNil(err)
	assert.Nil(new)
}
