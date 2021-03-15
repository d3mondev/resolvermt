package fastdns

import (
	"github.com/miekg/dns"
)

type resolverDNS struct {
	serverBalancer balancer
	parser         messageParser

	retryCount int
}

type balancer interface {
	Next() server
}

type messageParser interface {
	Parse(msg *dns.Msg) []Record
}

func newResolverDNS(retryCount int, serverBalancer balancer, parser messageParser) *resolverDNS {
	return &resolverDNS{
		serverBalancer: serverBalancer,
		parser:         parser,
		retryCount:     retryCount,
	}
}

func (s *resolverDNS) Resolve(query string, rrtype RRtype) []Record {
	records := []Record{}

	var err error
	var msg *dns.Msg

	// Send the request to a server, retrying on error
	for i := 0; i < s.retryCount; i++ {
		server := s.serverBalancer.Next()
		msg, _, err = server.Query(query, rrtype)

		if err == nil {
			break
		}
	}

	// Something went wrong with the request
	if err != nil {
		return records
	}

	// Request was successful, but a valid error such as NXDOMAIN occured
	if msg.Rcode != dns.RcodeSuccess {
		return records
	}

	records = s.parser.Parse(msg)

	return records
}
