package category

import (
	"context"
	"fmt"
	// "strings"
	"testing"

	// "github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCategoryCRUD(t *testing.T) {
	t.Skip()
	_ = assert.New(t)
	_ = context.Background()

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
			// todo: integration tests
		})
	}
}
