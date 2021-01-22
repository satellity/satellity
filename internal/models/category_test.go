package models

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCategoryCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer teardownTestContext(ctx)

	categoryCases := []struct {
		name        string
		alias       string
		description string
		position    int
		valid       bool
	}{
		{"", "alias", "description", 0, false},
		{"general", "", "description", 1, true},
		{"community", "    ", "description", 2, true},
		{"jobs", "Remote Jobs", "description", 3, true},
		{"Golang", "Golang", "", 4, true},
	}

	for _, tc := range categoryCases {
		t.Run(fmt.Sprintf("category name %s", tc.name), func(t *testing.T) {
			category, err := CreateCategory(ctx, tc.name, tc.alias, tc.description, 0)
			if !tc.valid {
				assert.NotNil(err)
				assert.Nil(category)
				return
			}

			assert.Nil(err)
			assert.NotNil(category)
			assert.Equal(tc.name, category.Name)
			if strings.TrimSpace(tc.alias) == "" {
				tc.alias = tc.name
			}
			assert.Equal(tc.alias, category.Alias)
			assert.Equal(tc.description, category.Description)
			new, err := ReadCategory(ctx, category.CategoryID)
			assert.Nil(err)
			assert.NotNil(new)
			new, err = ReadCategoryByIDOrName(ctx, category.Name)
			assert.Nil(err)
			assert.NotNil(new)
			assert.Equal(category.CategoryID, new.CategoryID)
			new, err = emitToCategory(ctx, category.CategoryID)
			assert.Nil(err)
			assert.NotNil(new)
			categories, err := ReadAllCategories(ctx)
			assert.Nil(err)
			assert.Len(categories, tc.position)
			new, err = ReadCategory(ctx, uuid.Must(uuid.NewV4()).String())
			assert.Nil(err)
			assert.Nil(new)
			new, err = UpdateCategory(ctx, uuid.Must(uuid.NewV4()).String(), "new"+category.Name, "new"+category.Alias, "new"+category.Description, 10)
			assert.Nil(err)
			assert.Nil(new)
			new, err = UpdateCategory(ctx, category.CategoryID, "", "", "", 10)
			assert.Nil(err)
			assert.NotNil(new)
			new, err = UpdateCategory(ctx, category.CategoryID, "new"+category.Name, "", "", 10)
			assert.Nil(err)
			assert.NotNil(new)
			assert.Equal("new"+tc.name, new.Name)
			assert.Equal(tc.alias, new.Alias)
			assert.Equal("", new.Description)
			new, err = UpdateCategory(ctx, category.CategoryID, "new"+category.Name, "new"+category.Alias, "", 10)
			assert.Nil(err)
			assert.NotNil(new)
			assert.Equal("new"+tc.name, new.Name)
			assert.Equal("new"+tc.alias, new.Alias)
			assert.Equal("", new.Description)
			new, err = UpdateCategory(ctx, category.CategoryID, "new"+category.Name, "new"+category.Alias, "new"+category.Description, 10)
			assert.Nil(err)
			assert.NotNil(new)
			assert.Equal("new"+tc.name, new.Name)
			assert.Equal("new"+tc.alias, new.Alias)
			assert.Equal("new"+tc.description, new.Description)
		})
	}
}
