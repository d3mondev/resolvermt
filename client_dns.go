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
	index := 0
	activeChan := make(chan bool, s.maxConcurrency)
	resultChan := make(chan []Record, s.maxConcurrency)
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
			<-activeChan
		}(query, rrtype, activeChan, resultChan)

		// Exit condition
		if index >= len(queries) {
			break
		}
	}

	// Process completed results
	for len(activeChan) > 0 || len(resultChan) > 0 {
		records = append(records, <-resultChan...)
	}

	// Work done
	return records
}
