package multidns

import (
	"time"

	"github.com/miekg/dns"
)

type request interface {
	Send(resolver string) (r *dns.Msg, rtt time.Duration, err error)
}

type requestDNS struct {
	c *dns.Client
	m *dns.Msg
}

type requestFactory func(query string, rrtype RRtype) request

func newRequestDNS(query string, rrtype RRtype) request {
	req := new(requestDNS)

	req.c = new(dns.Client)
	req.m = new(dns.Msg)
	req.m.SetQuestion(query+".", uint16(rrtype))

	return req
}

func (s *requestDNS) Send(resolver string) (*dns.Msg, time.Duration, error) {
	return s.c.Exchange(s.m, resolver)
}
