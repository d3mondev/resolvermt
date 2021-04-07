package resolvermt

import "sync/atomic"

type clientDNS struct {
	resolver       resolver
	maxConcurrency int

	queryCount int32
}

type resolver interface {
	Resolve(query string, rrtype RRtype) []Record
	Close()
}

func newClientDNS(resolver resolver, maxConcurrency int) *clientDNS {
	client := clientDNS{
		resolver:       resolver,
		maxConcurrency: maxConcurrency,
	}

	return &client
}

func (s *clientDNS) Resolve(queries []string, rrtype RRtype) []Record {
	records := []Record{}

	queryCount := len(queries)
	queryIndex := 0

	if queryCount == 0 {
		return records
	}

	queryChan := make(chan string, s.maxConcurrency)
	resultChan := make(chan []Record, s.maxConcurrency)

	received := 0

	// Start goroutines
	for i := 0; i < s.maxConcurrency; i++ {
		go func(queryChan chan string, resultChan chan []Record, rrtype RRtype) {
			for {
				query, open := <-queryChan

				if !open {
					return
				}

				results := s.resolver.Resolve(query, rrtype)
				resultChan <- results
			}
		}(queryChan, resultChan, rrtype)
	}

	// Send work to goroutines
	for {
		// Process completed results
		for i := len(resultChan); i > 0; i-- {
			records = append(records, <-resultChan...)
			received++
		}

		// Send the next query
		queryChan <- queries[queryIndex]
		queryIndex++

		// Exit condition
		if queryIndex >= queryCount {
			break
		}
	}

	// Process remaining results
	for received < queryCount {
		records = append(records, <-resultChan...)
		received++
	}

	// Work done
	close(queryChan)

	atomic.AddInt32(&s.queryCount, int32(queryCount))

	return records
}

func (s *clientDNS) QueryCount() int {
	return int(atomic.LoadInt32(&s.queryCount))
}

func (s *clientDNS) Close() {
	s.resolver.Close()
}
