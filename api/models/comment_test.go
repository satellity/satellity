package models

import (
	"context"
	"testing"
	"time"

	"github.com/go-pg/pg"
	"github.com/godiscourse/godiscourse/api/session"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCommentCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer session.Database(ctx).Close()
	defer teardownTestContext(ctx)

	user := createTestUser(ctx, "im.yuqlee@gmail.com", "username", "password")
	assert.NotNil(user)
	category, _ := CreateCategory(ctx, "name", "alias", "Description", 0)
	assert.NotNil(category)
	topic, _ := user.CreateTopic(ctx, "title", "body", category.CategoryID)
	assert.NotNil(topic)
	comment, err := user.CreateComment(ctx, uuid.Must(uuid.NewV4()).String(), "hello comment")
	assert.NotNil(err)
	assert.Nil(comment)
	comment, err = user.CreateComment(ctx, topic.TopicID, "hello comment")
	assert.Nil(err)
	assert.NotNil(comment)
	comment, err = readTestComment(ctx, comment.CommentID)
	assert.Nil(err)
	assert.NotNil(comment)
	topic, _ = ReadTopic(ctx, topic.TopicID)
	assert.NotNil(topic)
	assert.Equal(1, topic.CommentsCount)
	new, err := user.UpdateComment(ctx, comment.CommentID, "hello comment hello")
	assert.Nil(err)
	assert.NotNil(new)
	assert.Equal(comment.CommentID, new.CommentID)
	comments, err := topic.ReadComments(ctx, time.Time{})
	assert.Nil(err)
	assert.Len(comments, 1)
	comments, err = user.ReadComments(ctx, time.Time{})
	assert.Nil(err)
	assert.Len(comments, 1)
	comment, err = user.CreateComment(ctx, topic.TopicID, "hello comment 01")
	assert.Nil(err)
	assert.NotNil(comment)
	topic, _ = ReadTopic(ctx, topic.TopicID)
	assert.NotNil(topic)
	assert.Equal(2, topic.CommentsCount)
	err = user.DeleteComment(ctx, comment.CommentID)
	assert.Nil(err)
	topic, err = ReadTopic(ctx, topic.TopicID)
	assert.Nil(err)
	assert.NotNil(topic)
	assert.Equal(1, topic.CommentsCount)
	assert.Equal("title", topic.Title)

	user = createTestUser(ctx, "im.jadeydi@gmail.com", "usernamex", "password")
	assert.NotNil(user)
	comments, err = user.ReadComments(ctx, time.Time{})
	assert.Nil(err)
	assert.Len(comments, 0)
}

func readTestComment(ctx context.Context, id string) (*Comment, error) {
	comment := &Comment{CommentID: id}
	if err := session.Database(ctx).Model(comment).Column(commentColumns...).WherePK().Select(); err == pg.ErrNoRows {
		return nil, session.NotFoundError(ctx)
	} else if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return comment, nil
}
