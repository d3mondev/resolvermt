package multidns

import (
	"time"
)

// Client is used to send DNS requests to resolvers concurrently.
type Client interface {
	Resolve(domains []string, rrtype RRtype) []Record
}

// Record contains a DNS record entry returned by Client.Resolve.
type Record struct {
	Question string
	Type     RRtype
	Answer   string
}

// New returns a Client that will respect the retry count, queries per seconds
// and a maximum number of concurrent queries that can happen at the same time.
func New(resolvers []string, retryCount int, queriesPerSecond int, parallelCount int) Client {
	items := make([]server, len(resolvers))
	for i := range resolvers {
		items[i] = newRateLimitedServer(resolvers[i], queriesPerSecond)
	}
	roundRobinList := newRoundRobinList(items)

	parser := &msgParser{}
	resolver := newResolverDNS(retryCount, roundRobinList, parser)
	sleeper := &defaultSleeper{}

	return newClientDNS(resolver, sleeper, parallelCount)
}

// Resolver is used to resolve a DNS query and return a list of records.
// The process is responsible for selecting which DNS servers to use and
// rate-limiting.
type Resolver interface {
	Resolve(query string, rrtype RRtype) []Record
}

// Sleeper defines the behavior of the client when maximum concurrency is reached.
// Typically, the client will sleep for a number of milliseconds before processing
// more requests in order to let the Resolver finish a request.
type Sleeper interface {
	Sleep(t time.Duration)
}

type defaultSleeper struct{}

func (s *defaultSleeper) Sleep(t time.Duration) {
	time.Sleep(t)
}
