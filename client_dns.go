package fastdns

import (
	"time"
)

type clientDNS struct {
	resolver       Resolver
	maxConcurrency int
}

func newClientDNS(resolver Resolver, maxConcurrency int) *clientDNS {
	client := clientDNS{
		resolver:       resolver,
		maxConcurrency: maxConcurrency,
	}

	return &client
}

func (s *clientDNS) Resolve(queries []string, rrtype RRtype) []Record {
	// Limit the number of concurrent routines by using a channel
	activeChan := make(chan bool, s.maxConcurrency)

	// The result channel must have an extra element to prevent a deadlock
	// that happens when a goroutine tries to push its results in a channel that
	// is full while the main thread is waiting for the goroutines to finish
	resultChan := make(chan []Record, s.maxConcurrency+1)

	index := 0
	records := []Record{}

	for {
		// Process completed results
		for i := len(resultChan); i > 0; i-- {
			records = append(records, <-resultChan...)
		}

		// Get the next query
		query := queries[index]
		index++

		// Start a new goroutine
		activeChan <- true

		go func(query string, rrtype RRtype, activeChan chan bool, resultChan chan []Record) {
			records := s.resolver.Resolve(query, rrtype)
			resultChan <- records

			// Free up goroutine
			<-activeChan
		}(query, rrtype, activeChan, resultChan)

		// Exit condition
		if index >= len(queries) {
			break
		}
	}

	// Wait for routines to finish
	for len(activeChan) > 0 {
		time.Sleep(1 * time.Millisecond)
	}

	// Process completed results
	for i := len(resultChan); i > 0; i-- {
		records = append(records, <-resultChan...)
	}

	// Work done
	return records
}
