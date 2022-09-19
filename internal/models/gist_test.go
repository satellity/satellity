package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGistCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer teardownTestContext(ctx)

	source := testCreateSource(ctx)
	assert.NotNil(source)
	gist, err := CreateGist(ctx, "identity", "author", "gist title", "release", true, "link", "body", time.Now(), source)
	assert.Nil(err)
	assert.NotNil(gist)
	err = gist.Update(ctx, "gist title updated", "release", true)
	assert.Nil(err)
	gists, err := ReadGists(ctx, time.Now(), 64)
	assert.Nil(err)
	assert.Len(gists, 1)
	old, err := ReadGist(ctx, gist.GistID)
	assert.Nil(err)
	assert.NotNil(old)
	assert.Equal("gist title updated", old.Title)
}
