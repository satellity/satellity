package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSourceCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer teardownTestContext(ctx)

	source, err := CreateSource(ctx, "github", "https://github.com/AlphaWallet/alpha-wallet-ios/releases.atom", "logo")
	assert.Nil(err)
	assert.NotNil(source)

	source, err = CreateSource(ctx, "github", "https://github.com/AlphaWallet/alpha-wallet-ios/releases.atom", "logo")
	assert.Nil(err)
	assert.NotNil(source)

	sources, err := ReadSources(ctx)
	assert.Nil(err)
	assert.Len(sources, 1)

	err = source.Update(ctx, "jason", "host", "logo")
	assert.Nil(err)
	old, err := ReadSource(ctx, source.SourceID)
	assert.Nil(err)
	assert.NotNil(old)
	assert.Equal("jason", old.Author)
}
func testCreateSource(ctx context.Context) *Source {
	source, _ := CreateSource(ctx, "github", "https://github.com/AlphaWallet/alpha-wallet-ios/releases.atom", "logo")
	return source
}
