package multidns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNextResolver(t *testing.T) {
	testTable := []struct {
		name      string
		resolvers []string
		count     int
		want      string
	}{
		{name: "No Resolver", resolvers: []string{}, count: 1, want: ""},
		{name: "Nil Resolver", resolvers: nil, count: 1, want: ""},
		{name: "Single", resolvers: []string{"8.8.8.8:53"}, count: 1, want: "8.8.8.8:53"},
		{name: "Two", resolvers: []string{"8.8.8.8:53", "8.8.4.4:53"}, count: 2, want: "8.8.4.4:53"},
		{name: "Wrap", resolvers: []string{"8.8.8.8:53", "8.8.4.4:53"}, count: 3, want: "8.8.8.8:53"},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			resolverList := newResolverListRoundRobin(test.resolvers, 10)

			if len(test.resolvers) == 0 {
				assert.Nil(t, resolverList)
				return
			}

			var got resolver
			for i := 0; i < test.count; i++ {
				got = resolverList.GetNextResolver()
			}

			assert.Equal(t, test.want, got.Get())
		})
	}
}
