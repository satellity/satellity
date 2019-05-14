package engine

import (
	"context"
	"fmt"
	"godiscourse/internal/models"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCategoryCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	_, store, teardown := setup(t)
	defer teardown()

	categoryCases := []struct {
		models.CategoryRequest
		valid bool
	}{
		{models.CategoryRequest{Name: "", Alias: "alias", Description: "description", Position: 0}, false},
		{models.CategoryRequest{Name: "general", Alias: "", Description: "description", Position: 1}, true},
		{models.CategoryRequest{Name: "community", Alias: "    ", Description: "description", Position: 2}, true},
		{models.CategoryRequest{Name: "jobs", Alias: "Remote Jobs", Description: "description", Position: 3}, true},
		{models.CategoryRequest{Name: "Golang", Alias: "Golang", Description: "", Position: 4}, true},
	}

	for _, tc := range categoryCases {
		cr := tc.CategoryRequest
		t.Run(fmt.Sprintf("category name %s", cr.Name), func(t *testing.T) {
			category, err := store.CreateCategory(ctx, &cr)
			if !tc.valid {
				assert.NotNil(err)
				assert.Nil(category)
				return
			}

			assert.Nil(err)
			assert.NotNil(category)
			assert.Equal(cr.Name, category.Name)
			if strings.TrimSpace(cr.Alias) == "" {
				cr.Alias = cr.Name
			}
			assert.Equal(cr.Alias, category.Alias)
			assert.Equal(cr.Description, category.Description)
			new, err := store.GetCategoryByID(ctx, category.CategoryID)
			assert.Nil(err)
			assert.NotNil(new)
			assert.Equal(category.CategoryID, new.CategoryID)
			new, err = store.transmitToCategory(ctx, category.CategoryID)
			assert.Nil(err)
			assert.NotNil(new)
			categories, err := store.GetAllCategories(ctx)
			assert.Nil(err)
			assert.Len(categories, int(cr.Position))
			new, err = store.GetCategoryByID(ctx, uuid.Must(uuid.NewV4()).String())
			assert.Nil(err)
			assert.Nil(new)
			new, err = store.UpdateCategory(ctx, uuid.Must(uuid.NewV4()).String(), &cr)
			assert.NotNil(err)
			assert.Nil(new)
			copy := cr
			copy.Alias, copy.Name = "", ""
			new, err = store.UpdateCategory(ctx, category.CategoryID, &copy)
			assert.NotNil(err)
			assert.Nil(new)
			copy.Name = "new" + category.Name
			new, err = store.UpdateCategory(ctx, category.CategoryID, &copy)
			assert.Nil(err)
			assert.NotNil(new)
			assert.Equal(copy.Name, new.Name)
			assert.Equal(cr.Alias, new.Alias)
			assert.Equal(cr.Description, new.Description)
			copy.Alias = "new" + cr.Alias
			new, err = store.UpdateCategory(ctx, category.CategoryID, &copy)
			assert.Nil(err)
			assert.NotNil(new)
			assert.Equal(copy.Name, new.Name)
			assert.Equal(copy.Alias, new.Alias)
			assert.Equal(cr.Description, new.Description)
			copy.Description = "new" + cr.Description
			new, err = store.UpdateCategory(ctx, category.CategoryID, &copy)
			assert.Nil(err)
			assert.NotNil(new)
			assert.Equal(copy.Name, new.Name)
			assert.Equal(copy.Alias, new.Alias)
			assert.Equal(copy.Description, new.Description)
		})
	}
}
