package multidns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolverGet(t *testing.T) {
	testTable := []struct {
		name       string
		ipAddrPort string
		want       string
	}{
		{name: "IP Without Port", ipAddrPort: "8.8.8.8", want: "8.8.8.8:53"},
		{name: "IP With Port", ipAddrPort: "8.8.8.8:53", want: "8.8.8.8:53"},
		{name: "Empty IP", ipAddrPort: "", want: ":53"},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			resolver := newResolver(test.ipAddrPort, 10)
			got := resolver.Get()
			assert.Equal(t, test.want, got)
		})
	}
}
