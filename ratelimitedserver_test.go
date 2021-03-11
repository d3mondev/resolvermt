package multidns

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type spyLimiter struct {
	calls int32
}

func (s *spyLimiter) Take() time.Time {
	atomic.AddInt32(&s.calls, 1)

	return time.Time{}
}

func TestRateLimitedServerTake(t *testing.T) {
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
			limiter := &spyLimiter{}

			server := newRateLimitedServer(test.ipAddrPort, 10)
			server.limiter = limiter

			got := server.Take()
			assert.Equal(t, test.want, got)
			assert.Equal(t, 1, int(limiter.calls))
		})
	}
}
