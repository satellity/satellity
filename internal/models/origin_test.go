package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOriginCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer teardownTestContext(ctx)

	origin, err := CreateOrigin(ctx, "Twitter", "https://twitter.com/", "", "twitter")
	assert.Nil(err)
	assert.NotNil(origin)
	exist, err := ReadOrigin(ctx, origin.OriginID)
	assert.Nil(err)
	assert.NotNil(exist)
	assert.Equal(origin.OriginID, exist.OriginID)
	exist, err = CreateOrigin(ctx, "Twitter", "https://twitter.com/", "", "twitter")
	assert.Nil(err)
	assert.NotNil(exist)
	assert.Equal(origin.OriginID, exist.OriginID)
}
