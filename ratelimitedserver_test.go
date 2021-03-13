package fastdns

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type spyLimiter struct {
	calls int
}

func (s *spyLimiter) Take() time.Time {
	s.calls++

	return time.Time{}
}

func TestRateLimitedServerSend(t *testing.T) {
	limiter := &spyLimiter{}

	server := newRateLimitedServer("8.8.8.8:53", 10)
	server.limiter = limiter

	_, _, gotErr := server.Query("www.google.com", TypeA)

	assert.Equal(t, 1, limiter.calls)
	assert.Nil(t, gotErr)
}
