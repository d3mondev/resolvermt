package multidns

import (
	"strings"

	"go.uber.org/ratelimit"
)

type resolver interface {
	Get() string
}

func newResolver(ipAddrPort string, queriesPerSecond int) resolver {
	return newResolverRateLimited(ipAddrPort, queriesPerSecond)
}

type resolverRateLimited struct {
	ipAddrPort string
	limiter    ratelimit.Limiter
}

func newResolverRateLimited(ipAddrPort string, queriesPerSecond int) resolver {
	if !strings.Contains(ipAddrPort, ":") {
		ipAddrPort += ":53"
	}

	return &resolverRateLimited{
		ipAddrPort: ipAddrPort,
		limiter:    ratelimit.New(queriesPerSecond),
	}
}

// Get blocks to respect the rate limit and returns the resolver's ip:port string
func (s *resolverRateLimited) Get() string {
	s.limiter.Take()

	return s.ipAddrPort
}
