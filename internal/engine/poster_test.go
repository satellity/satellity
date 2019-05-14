package engine

import (
	"context"
	"testing"

	"godiscourse/internal/durable"
	"godiscourse/internal/models"
	"godiscourse/internal/user"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*durable.Database, *Store, func()) {
	db, teardown := setupTestContext()
	return db, NewStore(db), teardown
}

func TestUpdateUser(t *testing.T) {
	db, store, teardown := setup(t)
	defer teardown()

	ids := seedUsers(user.New(db), t)
	updatedCases := []struct {
		nickname  string
		biography string
	}{
		{"Dave", "I'm created from test"},
		{"Robert", "I'm updated from test!"},
		{"David", ""},
	}

	assert.Equal(t, len(ids), len(updatedCases))

	for i, id := range ids {
		err := store.UpdateUser(context.Background(), id, &models.UserInfo{
			Nickname:  updatedCases[i].nickname,
			Biography: updatedCases[i].biography,
		})

		assert.Nil(t, err)
	}
}
