package models

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSourceCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer teardownTestContext(ctx)

	sources, err := ReadSources(ctx)
	assert.Nil(err)
	assert.Len(sources, 0)

	source, err := CreateSource(ctx, "github", "https://github.com/AlphaWallet/alpha-wallet-ios/releases.atom", "logo", "locality")
	assert.Nil(err)
	assert.NotNil(source)
	assert.Equal("github.com", source.Host)

	source, err = CreateSource(ctx, "github", "https://github.com/AlphaWallet/alpha-wallet-ios/releases.atom", "logo", "locality")
	assert.Nil(err)
	assert.NotNil(source)

	sources, err = ReadSources(ctx)
	assert.Nil(err)
	assert.Len(sources, 1)

	err = source.Update(ctx, "jason", "host", "logo", 0, time.Now(), time.Now())
	assert.Nil(err)
	old, err := ReadSource(ctx, source.SourceID)
	assert.Nil(err)
	assert.NotNil(old)
	assert.Equal("jason", old.Author)

	err = source.Delete(ctx)
	assert.Nil(err)
	old, err = ReadSource(ctx, source.SourceID)
	assert.Nil(err)
	assert.Nil(old)
}

func testCreateSource(ctx context.Context) *Source {
	source, _ := CreateSource(ctx, "github", "https://github.com/AlphaWallet/alpha-wallet-ios/releases.atom", "logo", "locality")
	return source
}
