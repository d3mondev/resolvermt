package fastdns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRRTypeString(t *testing.T) {
	testTable := []struct {
		name   string
		rrtype RRtype
		want   string
	}{
		{name: "A", rrtype: TypeA, want: "A"},
		{name: "AAAA", rrtype: TypeAAAA, want: "AAAA"},
		{name: "CNAME", rrtype: TypeCNAME, want: "CNAME"},
		{name: "TXT", rrtype: TypeTXT, want: "TXT"},
		{name: "MX", rrtype: TypeMX, want: "MX"},
		{name: "NS", rrtype: TypeNS, want: "NS"},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			got := test.rrtype.String()
			assert.Equal(t, test.want, got)
		})
	}
}
