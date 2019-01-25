package models

import (
	"testing"

	"github.com/godiscourse/godiscourse/api/session"
	"github.com/stretchr/testify/assert"
)

func TestCategoryCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer session.Database(ctx).Close()
	defer teardownTestContext(ctx)

	category, err := CreateCategory(ctx, "name", "alias", "Description", 0)
	assert.NotNil(category)
	assert.Nil(err)
	new, err := ReadCategory(ctx, category.CategoryID)
	assert.NotNil(new)
	assert.Nil(err)
	assert.Equal("name", category.Name)
	assert.Equal("alias", category.Alias)
	assert.Equal(0, category.TopicsCount)
	assert.False(category.LastTopicID.Valid)
	category, err = UpdateCategory(ctx, category.CategoryID, "new name", "new alias", "new description", 0)
	assert.Nil(err)
	assert.NotNil(category)
	assert.Equal("new name", category.Name)
	assert.Equal("new alias", category.Alias)
	assert.Equal("new description", category.Description)
	category, err = CreateCategory(ctx, "name", "alias", "Description", 0)
	assert.NotNil(category)
	assert.Nil(err)
	assert.Equal(1, category.Position)
	categories, err := ReadCategories(ctx)
	assert.Nil(err)
	assert.Len(categories, 2)
}
