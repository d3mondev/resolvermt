package resolvermt

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
	Close()
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
	var i int
	for i = 0; i < s.retryCount; i++ {
		server := s.serverBalancer.Next()
		msg, _, err = server.Query(query, rrtype)

		if err != nil {
			continue
		}

		if msg.Rcode == dns.RcodeRefused || msg.Rcode == dns.RcodeServerFailure {
			continue
		}

		break
	}

	// Something went wrong with the request
	if i == s.retryCount {
		return records
	}

	// Request was successful, but a valid error such as NXDOMAIN occured
	if msg.Rcode != dns.RcodeSuccess {
		return records
	}

	records = s.parser.Parse(msg)

	return records
}

func (s *resolverDNS) Close() {
	s.serverBalancer.Close()
}
