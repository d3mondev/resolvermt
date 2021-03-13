package fastdns

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

type spyBalancer struct {
	calls  int32
	server server
}

func (s *spyBalancer) Next() server {
	s.calls++

	return s.server
}

type stubParser struct{}

func (s *stubParser) Parse(query string, msg *dns.Msg) []Record {
	return []Record{
		{
			Question: "test",
			Type:     TypeA,
			Answer:   "127.0.0.1",
		},
	}
}

type fakeServer struct {
	errors int
	msg    *dns.Msg
}

func (s *fakeServer) Query(query string, ttype RRtype) (*dns.Msg, time.Duration, error) {
	if s.errors > 0 {
		s.errors--
		return nil, time.Duration(0), errors.New("error")
	}

	return s.msg, time.Duration(0), nil
}

func TestResolve(t *testing.T) {
	stubRecords := []Record{
		{"test", TypeA, "127.0.0.1"},
	}

	stubA := &dns.A{}
	stubA.Hdr.Rrtype = dns.TypeA
	stubA.A = net.ParseIP("127.0.0.1")

	stubMsg := &dns.Msg{}
	stubMsg.Answer = []dns.RR{stubA}

	stubMsgErr := &dns.Msg{}
	stubMsgErr.Rcode = dns.RcodeServerFailure

	testTable := []struct {
		name    string
		retries int
		errors  int
		msg     *dns.Msg
		want    []Record
	}{
		{name: "Single Answer", retries: 3, errors: 0, msg: stubMsg, want: stubRecords},
		{name: "Retry", retries: 3, errors: 2, msg: stubMsg, want: stubRecords},
		{name: "Max Retry", retries: 1, errors: 1, msg: stubMsg, want: []Record{}},
		{name: "DNS Error", retries: 3, errors: 0, msg: stubMsgErr, want: []Record{}},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			stubMessageParser := &stubParser{}
			fakeServer := &fakeServer{errors: test.errors, msg: test.msg}
			spyBalancer := &spyBalancer{server: fakeServer}

			resolver := newResolverDNS(test.retries, spyBalancer, stubMessageParser)

			got := resolver.Resolve("test", TypeA)

			assert.EqualValues(t, test.want, got, test.name)

			wantedBalancerCalls := test.errors
			if test.retries > test.errors {
				wantedBalancerCalls++
			}

			assert.EqualValues(t, wantedBalancerCalls, spyBalancer.calls)
		})
	}
}
