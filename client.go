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
	doer := newResolver(retryCount, realNewSender, roundRobinList, parser)
	sleeper := &realSleeper{}

	return newClientDNS(doer, sleeper, parallelCount)
}

type doer interface {
	Resolve(query string, rrtype RRtype, channel chan []Record)
}

type sleeper interface {
	Sleep(t time.Duration)
}

type realSleeper struct{}

func (s *realSleeper) Sleep(t time.Duration) {
	time.Sleep(t)
}

func realNewSender(query string, rrtype RRtype) sender {
	return newRequestDNS(query, rrtype)
}
