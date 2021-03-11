package multidns

import (
	"strings"

	"go.uber.org/ratelimit"
)

type rateLimitedServer struct {
	ipAddrPort string
	limiter    ratelimit.Limiter
}

type server interface {
	// Take returns a server's IP:Port string, and may be blocking in order
	// to respect the rate limit.
	Take() string
}

func newRateLimitedServerList(ipAddrPort []string, queriesPerSecond int) []server {
	list := make([]server, len(ipAddrPort))

	for i := range ipAddrPort {
		list[i] = newRateLimitedServer(ipAddrPort[i], queriesPerSecond)
	}

	return list
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

func (s *rateLimitedServer) Take() string {
	s.limiter.Take()

	return s.ipAddrPort
}
