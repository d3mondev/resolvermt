package multidns

import (
	"time"
)

// Client interface
type Client interface {
	Resolve(domains []string, rrtype RRtype) []Record
}

// New creates a new concrete instance of Client
func New(resolvers []string, retryCount int, queriesPerSecond int, parallelCount int) Client {
	resolverList := newResolverListRoundRobin(resolvers, queriesPerSecond)
	parser := &msgParser{}
	requestFactory := newRequestDNS

	return newClientDNS(resolverList, requestFactory, parser, retryCount, parallelCount)
}

type clientDNS struct {
	resolvers      resolverList
	requestFactory requestFactory
	parser         parser

	retryCount  int
	workerCount int
}

func newClientDNS(resolvers resolverList, requestFactory requestFactory, parser parser, retryCount int, parallelCount int) Client {
	client := clientDNS{
		resolvers:      resolvers,
		requestFactory: requestFactory,
		parser:         parser,

		retryCount:  retryCount,
		workerCount: parallelCount,
	}

	return &client
}

func (s *clientDNS) Resolve(queries []string, rrtype RRtype) []Record {
	index := 0
	channel := make(chan []Record, s.workerCount)
	activeRoutines := 0
	records := []Record{}

	for {
		// Free up completed goroutines
		for i := len(channel); i > 0; i-- {
			records = append(records, <-channel...)
			activeRoutines--
		}

		// Wait if too many routines are in flight
		if activeRoutines >= s.workerCount {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		// Get the next query
		query := queries[index]
		index++

		// Start a new goroutine
		activeRoutines++
		go work(query, rrtype, channel, s.retryCount, s.requestFactory, s.resolvers, s.parser)

		// Exit condition
		if index >= len(queries) {
			break
		}
	}

	// Wait for all goroutines to finish
	for i := activeRoutines; i > 0; i-- {
		records = append(records, <-channel...)
	}

	// Work done
	return records
}
