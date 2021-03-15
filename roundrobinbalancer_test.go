package fastdns

import (
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

type stubServer struct {
	value string
}

func (s *stubServer) Query(query string, rrtype RRtype) (*dns.Msg, time.Duration, error) {
	return &dns.Msg{}, time.Duration(0), nil
}

func TestRoundRobinBalancerNext(t *testing.T) {
	stubServerA := &stubServer{"8.8.8.8:53"}
	stubServerB := &stubServer{"8.8.4.4:53"}

	testTable := []struct {
		name        string
		haveServers []server
		haveCount   int
		want        server
	}{
		{name: "No Items", haveServers: []server{}, haveCount: 1, want: nil},
		{name: "Nil Resolver", haveServers: nil, haveCount: 1, want: nil},
		{name: "Single", haveServers: []server{stubServerA}, haveCount: 1, want: stubServerA},
		{name: "Second", haveServers: []server{stubServerA, stubServerB}, haveCount: 2, want: stubServerB},
		{name: "Wrap Around", haveServers: []server{stubServerA, stubServerB}, haveCount: 3, want: stubServerA},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			list := newRoundRobinBalancer(test.haveServers)

			var got server
			for i := 0; i < test.haveCount; i++ {
				got = list.Next()
			}

			assert.Equal(t, test.want, got)
		})
	}
}
