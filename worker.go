package multidns

import "github.com/miekg/dns"

type worker func(query string, rrtype RRtype, channel chan []Record, retryCount int, requestFactory requestFactory, resolvers resolverList, parser parser)

func work(query string, rrtype RRtype, channel chan []Record, retryCount int, requestFactory requestFactory, resolvers resolverList, parser parser) {
	request := requestFactory(query, rrtype)

	var err error
	var msg *dns.Msg

	// Send the request to a resolver, retrying on error
	for i := 0; i < retryCount; i++ {
		resolver := resolvers.GetNextResolver()
		msg, _, err = request.Send(resolver.Get())

		if err == nil {
			break
		}
	}

	// Something went wrong with the request
	if err != nil {
		channel <- []Record{}
		return
	}

	// Request was successful, but a valid error such as NXDOMAIN occured
	if msg.Rcode != dns.RcodeSuccess {
		channel <- []Record{}
		return
	}

	records := parser.Parse(query, msg)

	channel <- records
}
