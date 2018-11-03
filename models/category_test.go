package models

import (
	"testing"

	"github.com/godiscourse/godiscourse/session"
	"github.com/stretchr/testify/assert"
)

func TestCategoryCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer session.Database(ctx).Close()
	defer teardownTestContext(ctx)

	category, err := CreateCategory(ctx, "name", "Description")
	assert.NotNil(category)
	assert.Nil(err)
	new, err := ReadCategory(ctx, category.CategoryID)
	assert.NotNil(new)
	assert.Nil(err)
	assert.Equal("name", category.Name)
	assert.Equal(0, category.TopicsCount)
	assert.False(category.LastTopicID.Valid)
}
