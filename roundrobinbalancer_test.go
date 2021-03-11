package fastdns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubServer struct {
	value string
}

func (s *stubServer) Take() string {
	return s.value
}

func TestRoundRobinListNext(t *testing.T) {
	stubServerA := &stubServer{"8.8.8.8:53"}
	stubServerB := &stubServer{"8.8.4.4:53"}

	testTable := []struct {
		name    string
		servers []server
		count   int
		want    *stubServer
	}{
		{name: "No Items", servers: []server{}, count: 1, want: nil},
		{name: "Nil Resolver", servers: nil, count: 1, want: nil},
		{name: "Single", servers: []server{stubServerA}, count: 1, want: stubServerA},
		{name: "Second", servers: []server{stubServerA, stubServerB}, count: 2, want: stubServerB},
		{name: "Wrap Around", servers: []server{stubServerA, stubServerB}, count: 3, want: stubServerA},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			list := newRoundRobinBalancer(test.servers)

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
