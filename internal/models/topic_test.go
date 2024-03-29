package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTopicCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer teardownTestContext(ctx)

	user := createTestUser(ctx, "im.yuqlee@gmail.com", "password")
	assert.NotNil(user)
	category, _ := CreateCategory(ctx, "name", "alias", "Description", 0)
	assert.NotNil(category)

	topicCases := []struct {
		title          string
		body           string
		categoryID     string
		topicsCount    int64
		commentsCount  int64
		bookmarksCount int64
		likesCount     int64
		draft          bool
		hasDraft       bool
		valid          bool
	}{
		{"", "body", category.CategoryID, 0, 0, 0, 0, false, false, false},
		{"title", "body", uuid.Must(uuid.NewV4()).String(), 0, 0, 0, 0, false, false, false},
		{"title", "body", category.CategoryID, 1, 0, 0, 0, false, false, true},
		{"title2", "body", category.CategoryID, 2, 0, 0, 0, false, true, true},
	}

	for _, tc := range topicCases {
		t.Run(fmt.Sprintf("topic title %s", tc.title), func(t *testing.T) {
			if !tc.valid {
				topic, err := user.CreateTopic(ctx, tc.title, tc.body, TopicTypePost, tc.categoryID, tc.draft)
				assert.NotNil(err)
				assert.Nil(topic)
				return
			}

			topic, err := user.CreateTopic(ctx, tc.title, tc.body, TopicTypePost, category.CategoryID, tc.draft)
			assert.Nil(err)
			assert.NotNil(topic)
			time.Sleep(100 * time.Millisecond)
			category, _ = ReadCategory(ctx, category.CategoryID)
			assert.NotNil(category)
			//assert.Equal(topic.TopicID, *category.LastTopicID)
			assert.Equal(tc.topicsCount, category.TopicsCount)
			topic, err = ReadTopic(ctx, topic.TopicID)
			assert.Nil(err)
			assert.NotNil(topic)
			assert.Equal(tc.title, topic.Title)
			assert.Equal(tc.body, topic.Body)
			assert.Equal(tc.bookmarksCount, topic.BookmarksCount)
			assert.Equal(tc.likesCount, topic.LikesCount)
			existing, err := ReadTopic(ctx, uuid.Must(uuid.NewV4()).String())
			assert.Nil(err)
			assert.Nil(existing)
			topics, err := ReadTopics(ctx, time.Time{}, nil, nil)
			assert.Nil(err)
			assert.Len(topics, int(tc.topicsCount))
			topics, err = ReadTopics(ctx, time.Time{}, nil, user)
			assert.Nil(err)
			assert.Len(topics, int(tc.topicsCount))
			topics, err = ReadTopics(ctx, time.Time{}, category, nil)
			assert.Nil(err)
			assert.Len(topics, int(tc.topicsCount))

			topic, err = user.UpdateTopic(ctx, topic.TopicID, "hell", "orld", TopicTypePost, "", tc.draft)
			assert.Nil(err)
			assert.NotNil(topic)
			assert.Equal("hell", topic.Title)
			assert.Equal("orld", topic.Body)
			topic, err = user.UpdateTopic(ctx, topic.TopicID, topic.Title, "orld orld", TopicTypePost, topic.CategoryID, tc.draft)
			assert.Nil(err)
			assert.NotNil(topic)
			assert.Equal("hell", topic.Title)
			assert.Equal("orld orld", topic.Body)
			existing, err = user.UpdateTopic(ctx, uuid.Must(uuid.NewV4()).String(), "hell", "orld", TopicTypePost, "", tc.draft)
			assert.Nil(err)
			assert.Nil(existing)
			u := &User{UserID: uuid.Must(uuid.NewV4()).String()}
			existing, err = u.UpdateTopic(ctx, topic.TopicID, "hell", "orld", TopicTypePost, "", tc.draft)
			assert.NotNil(err)
			assert.Nil(existing)

			if !tc.hasDraft {
				topic, err = user.DraftTopic(ctx)
				assert.Nil(err)
				assert.Nil(topic)
				topic, err = user.CreateTopic(ctx, tc.title, tc.body, TopicTypePost, category.CategoryID, true)
				assert.Nil(err)
				assert.NotNil(topic)
				topic, err = user.DraftTopic(ctx)
				assert.Nil(err)
				assert.NotNil(topic)
			}

			if tc.hasDraft {
				topic, err = user.DraftTopic(ctx)
				assert.Nil(err)
				assert.NotNil(topic)
				topic, err = user.CreateTopic(ctx, tc.title, tc.body, TopicTypePost, category.CategoryID, true)
				assert.NotNil(err)
				assert.Nil(topic)
			}
		})
	}

	user = createTestUser(ctx, "im.jadeydi@gmail.com", "password")
	assert.NotNil(user)
	category, _ = CreateCategory(ctx, "new name", "new alias", "New Description", 2)
	assert.NotNil(category)
	topicCases = []struct {
		title          string
		body           string
		categoryID     string
		topicsCount    int64
		commentsCount  int64
		bookmarksCount int64
		likesCount     int64
		draft          bool
		hasDraft       bool
		valid          bool
	}{
		{"title", "body", category.CategoryID, 1, 0, 0, 0, false, false, true},
		{"title2", "body", category.CategoryID, 2, 0, 0, 0, false, false, true},
	}

	for _, tc := range topicCases {
		t.Run(fmt.Sprintf("topic title %s", tc.title), func(t *testing.T) {
			topic, err := user.CreateTopic(ctx, tc.title, tc.body, TopicTypePost, category.CategoryID, tc.draft)
			assert.Nil(err)
			assert.NotNil(topic)
			topics, err := ReadTopics(ctx, time.Time{}, nil, nil)
			assert.Nil(err)
			assert.Len(topics, int(tc.topicsCount+2))
			topics, err = ReadTopics(ctx, time.Time{}, nil, user)
			assert.Nil(err)
			assert.Len(topics, int(tc.topicsCount))
			topics, err = ReadTopics(ctx, time.Time{}, category, nil)
			assert.Nil(err)
			assert.Len(topics, int(tc.topicsCount))

			tu, err := readTopicUser(ctx, topic.TopicID, user.UserID)
			assert.Nil(err)
			assert.Nil(tu)
			topic, err = topic.ActiondBy(ctx, user, TopicUserActionLiked, true)
			assert.Nil(err)
			assert.True(topic.IsLikedBy)
			assert.False(topic.IsBookmarkedBy)
			assert.Equal(int64(1), topic.LikesCount)
			tu, err = readTopicUser(ctx, topic.TopicID, user.UserID)
			assert.Nil(err)
			assert.NotNil(tu)
			assert.True(tu.LikedAt.Valid)
			assert.False(tu.BookmarkedAt.Valid)
			topic, err = topic.ActiondBy(ctx, user, TopicUserActionBookmarked, true)
			assert.Nil(err)
			assert.True(topic.IsLikedBy)
			assert.True(topic.IsBookmarkedBy)
			assert.Equal(int64(1), topic.BookmarksCount)
			tu, err = readTopicUser(ctx, topic.TopicID, user.UserID)
			assert.Nil(err)
			assert.NotNil(tu)
			assert.True(tu.LikedAt.Valid)
			assert.True(tu.BookmarkedAt.Valid)
			topic, err = ReadTopic(ctx, topic.TopicID)
			assert.Nil(err)
			assert.NotNil(topic)
			assert.Equal(int64(1), topic.LikesCount)
			assert.Equal(int64(1), topic.BookmarksCount)
			topic, err = topic.ActiondBy(ctx, user, TopicUserActionLiked, false)
			assert.Nil(err)
			assert.False(topic.IsLikedBy)
			assert.True(topic.IsBookmarkedBy)
			assert.Equal(int64(0), topic.LikesCount)
			assert.Equal(int64(1), topic.BookmarksCount)
			tu, err = readTopicUser(ctx, topic.TopicID, user.UserID)
			assert.Nil(err)
			assert.NotNil(tu)
			assert.False(tu.LikedAt.Valid)
			assert.True(tu.BookmarkedAt.Valid)
			topic, err = topic.ActiondBy(ctx, user, TopicUserActionBookmarked, false)
			assert.Nil(err)
			assert.False(topic.IsLikedBy)
			assert.False(topic.IsBookmarkedBy)
			assert.Equal(int64(0), topic.LikesCount)
			assert.Equal(int64(0), topic.BookmarksCount)
			tu, err = readTopicUser(ctx, topic.TopicID, user.UserID)
			assert.Nil(err)
			assert.NotNil(tu)
			assert.False(tu.LikedAt.Valid)
			assert.False(tu.BookmarkedAt.Valid)
		})
	}
}
