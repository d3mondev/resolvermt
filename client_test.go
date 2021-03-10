package multidns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	resolvers := []string{"8.8.8.8"}
	client := New(resolvers, 3, 10, 5)
	assert.NotNil(t, client)
}
