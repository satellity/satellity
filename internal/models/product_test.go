package models

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestProductCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer teardownTestContext(ctx)

	user := createTestUser(ctx, "im.yuqlee@gmail.com", "username", "password")
	assert.NotNil(user)

	productCases := []struct {
		name   string
		body   string
		cover  string
		source string
		tags   pq.StringArray
		valid  bool
	}{
		{"", "body", "", "", []string{}, false},
		{"name", "body", "cover", "source", []string{"Big", "Small"}, true},
	}

	for _, pc := range productCases {
		t.Run(fmt.Sprintf("product title %s", pc.name), func(t *testing.T) {
			if !pc.valid {
				product, err := user.CreateProduct(ctx, pc.name, pc.body, pc.cover, pc.source, pc.tags)
				assert.NotNil(err)
				assert.Nil(product)
				return
			}

			product, err := user.CreateProduct(ctx, pc.name, pc.body, pc.cover, pc.source, pc.tags)
			assert.Nil(err)
			assert.NotNil(product)
			product, err = user.UpdateProduct(ctx, product.ProductID, "new"+pc.name, "new"+pc.body, "new"+pc.cover, "new"+pc.source, pc.tags)
			assert.Nil(err)
			assert.NotNil(product)
			product, err = FindProduct(ctx, product.ProductID)
			assert.Nil(err)
			assert.NotNil(product)
			assert.Equal("new"+pc.name, product.Name)
			assert.Equal("new"+pc.body, product.Body)

			products, err := FindProducts(ctx)
			assert.Nil(err)
			assert.Len(products, 1)
			products, err = RelatedProducts(ctx, uuid.Nil.String())
			assert.Nil(err)
			assert.Len(products, 1)

			products, err = SearchProducts(ctx, "Big")
			assert.Nil(err)
			assert.Len(products, 1)
			products, err = SearchProducts(ctx, "Big,Small")
			assert.Nil(err)
			assert.Len(products, 1)
			products, err = SearchProducts(ctx, "small")
			assert.Nil(err)
			assert.Len(products, 0)
		})
	}
}
