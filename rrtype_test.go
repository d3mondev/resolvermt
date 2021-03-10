package multidns_test

import (
	"testing"

	"github.com/d3mondev/multidns"
	"github.com/stretchr/testify/assert"
)

func TestRRTypeString(t *testing.T) {
	testTable := []struct {
		name   string
		rrtype multidns.RRtype
		want   string
	}{
		{name: "A", rrtype: multidns.TypeA, want: "A"},
		{name: "AAAA", rrtype: multidns.TypeAAAA, want: "AAAA"},
		{name: "CNAME", rrtype: multidns.TypeCNAME, want: "CNAME"},
		{name: "TXT", rrtype: multidns.TypeTXT, want: "TXT"},
		{name: "MX", rrtype: multidns.TypeMX, want: "MX"},
		{name: "NS", rrtype: multidns.TypeNS, want: "NS"},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			got := test.rrtype.String()
			assert.Equal(t, test.want, got)
		})
	}
}
