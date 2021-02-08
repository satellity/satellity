package models

import (
	"context"
	"fmt"
	"satellity/internal/session"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

func TestCommentCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer teardownTestContext(ctx)

	user := createTestUser(ctx, "im.yuqlee@gmail.com", "username", "password")
	assert.NotNil(user)
	category, _ := CreateCategory(ctx, "name", "alias", "Description", 0)
	assert.NotNil(category)
	topic, err := user.CreateTopic(ctx, "title", "body", TopicTypePost, category.CategoryID, false)
	assert.Nil(err)
	assert.NotNil(topic)

	commentCases := []struct {
		topicID string
		body    string
		valid   bool
	}{
		{topic.TopicID, "", false},
		{topic.TopicID, "      ", false},
		{uuid.Must(uuid.NewV4()).String(), "comment body", false},
		{topic.TopicID, "comment body", true},
	}

	for _, tc := range commentCases {
		t.Run(fmt.Sprintf("comment body %s", tc.body), func(t *testing.T) {
			if !tc.valid {
				comment, err := user.CreateComment(ctx, tc.topicID, tc.body)
				assert.NotNil(err)
				assert.Nil(comment)
				return
			}

			comment, err := user.CreateComment(ctx, tc.topicID, tc.body)
			assert.Nil(err)
			assert.NotNil(comment)
			assert.Equal(tc.body, comment.Body)
			new, err := readTestComment(ctx, comment.CommentID)
			assert.Nil(err)
			assert.NotNil(new)
			new, err = readTestComment(ctx, uuid.Must(uuid.NewV4()).String())
			assert.Nil(err)
			assert.Nil(new)
			new, err = user.UpdateComment(ctx, uuid.Must(uuid.NewV4()).String(), "comment body")
			assert.NotNil(err)
			assert.Nil(new)
			new, err = user.UpdateComment(ctx, comment.CommentID, "    ")
			assert.NotNil(err)
			assert.Nil(new)
			new, err = user.UpdateComment(ctx, comment.CommentID, "new comment body")
			assert.Nil(err)
			assert.NotNil(new)
			assert.Equal("new comment body", new.Body)
			comments, err := topic.ReadComments(ctx, time.Time{})
			assert.Nil(err)
			assert.Len(comments, 1)
			comments, err = user.ReadComments(ctx, time.Time{})
			assert.Nil(err)
			assert.Len(comments, 1)
			topic, _ = ReadTopic(ctx, topic.TopicID)
			assert.NotNil(topic)
			assert.Equal(int64(1), topic.CommentsCount)
			err = user.DeleteComment(ctx, comment.CommentID)
			assert.Nil(err)
			topic, err = ReadTopic(ctx, topic.TopicID)
			assert.Nil(err)
			assert.NotNil(topic)
			assert.Equal(int64(0), topic.CommentsCount)
			comments, err = topic.ReadComments(ctx, time.Time{})
			assert.Nil(err)
			assert.Len(comments, 0)
			comments, err = user.ReadComments(ctx, time.Time{})
			assert.Nil(err)
			assert.Len(comments, 0)
			new, err = readTestComment(ctx, comment.CommentID)
			assert.Nil(err)
			assert.Nil(new)
		})
	}
}

func readTestComment(ctx context.Context, id string) (*Comment, error) {
	var comment *Comment
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		comment, err = findComment(ctx, tx, id)
		return err
	})
	return comment, err
}
