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
	queryChan := make(chan []string, s.maxConcurrency)
	resultChan := make(chan []Record, s.maxConcurrency)

	queryCount := len(queries)
	queryIndex := 0

	batchSent := 0
	batchReceived := 0

	records := []Record{}

	// Start goroutines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < s.maxConcurrency; i++ {
		go func(ctx context.Context, queryChan chan []string, resultChan chan []Record, rrtype RRtype) {
			for {
				select {
				case <-ctx.Done():
					return
				case batch := <-queryChan:
					var results []Record

					for _, query := range batch {
						results = append(results, s.resolver.Resolve(query, rrtype)...)
					}

					resultChan <- results
				}
			}
		}(ctx, queryChan, resultChan, rrtype)
	}

	// Send work to goroutines
	for {
		// Process completed results
		for i := len(resultChan); i > 0; i-- {
			records = append(records, <-resultChan...)
			batchReceived++
		}

		// Send the next queries
		endIndex := queryIndex + batchSize(queryCount-queryIndex, s.maxConcurrency)
		slice := queries[queryIndex:endIndex]
		queryIndex = endIndex

		queryChan <- slice
		batchSent++

		// Exit condition
		if queryIndex >= queryCount {
			break
		}
	}

	// Process remaining results
	for batchReceived < batchSent {
		records = append(records, <-resultChan...)
		batchReceived++
	}

	// Work done
	return records
}

func batchSize(count int, threads int) int {
	batchsize := max(count/threads/2, 1)
	batchsize = min(batchsize, 100)

	return batchsize
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
