package fastdns

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

	queryCount := len(queries)
	queryIndex := 0

	received := 0

	records := []Record{}

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

	return records
}

func (s *clientDNS) Close() {
	s.resolver.Close()
}
