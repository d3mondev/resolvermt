package multidns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockResolver struct {
	value string
}

func (s *mockResolver) Take() string {
	return s.value
}

func TestRoundRobinListNext(t *testing.T) {
	mockResolverA := &mockResolver{"8.8.8.8:53"}
	mockResolverB := &mockResolver{"8.8.4.4:53"}

	testTable := []struct {
		name    string
		servers []server
		count   int
		want    *mockResolver
	}{
		{name: "No Items", servers: []server{}, count: 1, want: nil},
		{name: "Nil Resolver", servers: nil, count: 1, want: nil},
		{name: "Single", servers: []server{mockResolverA}, count: 1, want: mockResolverA},
		{name: "Second", servers: []server{mockResolverA, mockResolverB}, count: 2, want: mockResolverB},
		{name: "Wrap Around", servers: []server{mockResolverA, mockResolverB}, count: 3, want: mockResolverA},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			list := newRoundRobinList(test.servers)

			var got server
			for i := 0; i < test.count; i++ {
				got = list.Next()
			}

			if test.want == nil {
				assert.Nil(t, got)
			} else {
				assert.Equal(t, test.want, got)
			}
		})
	}
}
