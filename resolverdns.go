package multidns

import (
	"time"

	"github.com/miekg/dns"
)

type resolverDNS struct {
	serverBalancer balancer
	newSender      newSenderFunc
	parser         messageParser

	retryCount int
}

type balancer interface {
	Next() server
}

type sender interface {
	Send(resolver string) (r *dns.Msg, rtt time.Duration, err error)
}

type messageParser interface {
	Parse(query string, msg *dns.Msg) []Record
}

type newSenderFunc func(query string, rrtype RRtype) sender

func newResolverDNS(retryCount int, newSender newSenderFunc, serverBalancer balancer, parser messageParser) *resolverDNS {
	return &resolverDNS{
		serverBalancer: serverBalancer,
		newSender:      newSender,
		parser:         parser,
		retryCount:     retryCount,
	}
}

func (s *resolverDNS) Resolve(query string, rrtype RRtype) []Record {
	records := []Record{}
	sender := s.newSender(query, rrtype)

	var err error
	var msg *dns.Msg

	// Send the request to a server, retrying on error
	for i := 0; i < s.retryCount; i++ {
		server := s.serverBalancer.Next()
		msg, _, err = sender.Send(server.Take())

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

	records = s.parser.Parse(query, msg)

	return records
}
