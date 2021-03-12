package fastdns

import (
	"context"
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
	queryChan := make(chan string, s.maxConcurrency)
	resultChan := make(chan []Record, s.maxConcurrency)

	queryIndex := 0
	resultCount := 0

	records := []Record{}

	// Start goroutines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < s.maxConcurrency; i++ {
		go func(ctx context.Context, queryChan chan string, resultChan chan []Record, rrtype RRtype) {
			for {
				select {
				case <-ctx.Done():
					return
				case query := <-queryChan:
					records := s.resolver.Resolve(query, rrtype)
					resultChan <- records
				}
			}
		}(ctx, queryChan, resultChan, rrtype)
	}

	// Send work to goroutines
	for {
		// Process completed results
		for i := len(resultChan); i > 0; i-- {
			records = append(records, <-resultChan...)
			resultCount++
		}

		// Send the next query
		queryChan <- queries[queryIndex]
		queryIndex++

		// Exit condition
		if queryIndex >= len(queries) {
			break
		}
	}

	// Process remaining results
	for resultCount < len(queries) {
		records = append(records, <-resultChan...)
		resultCount++
	}

	// Work done
	return records
}
