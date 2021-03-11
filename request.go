package fastdns

import (
	"time"

	"github.com/miekg/dns"
)

type requestDNS struct {
	c *dns.Client
	m *dns.Msg
}

func newRequestDNS(query string, rrtype RRtype) *requestDNS {
	req := requestDNS{}
	req.c = new(dns.Client)
	req.m = new(dns.Msg)
	req.m.SetQuestion(query+".", uint16(rrtype))

	return &req
}

func (s *requestDNS) Send(server string) (*dns.Msg, time.Duration, error) {
	return s.c.Exchange(s.m, server)
}
