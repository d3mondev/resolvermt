package multidns

import (
	"strings"

	"go.uber.org/ratelimit"
)

type rateLimitedServer struct {
	ipAddrPort string
	limiter    ratelimit.Limiter
}

func newRateLimitedServer(ipAddrPort string, queriesPerSecond int) *rateLimitedServer {
	if !strings.Contains(ipAddrPort, ":") {
		ipAddrPort += ":53"
	}

	return &rateLimitedServer{
		ipAddrPort: ipAddrPort,
		limiter:    ratelimit.New(queriesPerSecond),
	}
}

// Take returns a server's IP:Port string, and may be blocking in order
// to respect the rate limit.
func (s *rateLimitedServer) Take() string {
	s.limiter.Take()

	return s.ipAddrPort
}
