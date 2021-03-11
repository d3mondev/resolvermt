package multidns

import "time"

type clientDNS struct {
	resolver resolver
	sleeper  sleeper

	workerCount int
}

func newClientDNS(resolver resolver, sleeper sleeper, parallelCount int) *clientDNS {
	client := clientDNS{
		resolver: resolver,
		sleeper:  sleeper,

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
			s.sleeper.Sleep(10 * time.Millisecond)
			continue
		}

		// Get the next query
		query := queries[index]
		index++

		// Start a new goroutine
		activeRoutines++
		go func(query string, rrtype RRtype, channel chan []Record) {
			records := s.resolver.Resolve(query, rrtype)
			channel <- records
		}(query, rrtype, channel)

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
