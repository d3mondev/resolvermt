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

func (s *stubParser) Parse(msg *dns.Msg) []Record {
	return []Record{
		{
			Question: msg.Question[0].Name,
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
	stubMsg.Question = []dns.Question{{Name: "test"}}
	stubMsg.Answer = []dns.RR{stubA}

	stubMsgErr := &dns.Msg{}
	stubMsgErr.Rcode = dns.RcodeServerFailure

	testTable := []struct {
		name        string
		haveRetries int
		haveErrors  int
		haveMsg     *dns.Msg
		want        []Record
	}{
		{name: "Single Answer", haveRetries: 3, haveErrors: 0, haveMsg: stubMsg, want: stubRecords},
		{name: "Retry", haveRetries: 3, haveErrors: 2, haveMsg: stubMsg, want: stubRecords},
		{name: "Max Retry", haveRetries: 1, haveErrors: 1, haveMsg: stubMsg, want: []Record{}},
		{name: "DNS Error", haveRetries: 3, haveErrors: 0, haveMsg: stubMsgErr, want: []Record{}},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			stubMessageParser := &stubParser{}
			fakeServer := &fakeServer{errors: test.haveErrors, msg: test.haveMsg}
			spyBalancer := &spyBalancer{server: fakeServer}

			resolver := newResolverDNS(test.haveRetries, spyBalancer, stubMessageParser)

			got := resolver.Resolve("test", TypeA)

			assert.EqualValues(t, test.want, got, test.name)

			wantedBalancerCalls := test.haveErrors
			if test.haveRetries > test.haveErrors {
				wantedBalancerCalls++
			}

			assert.EqualValues(t, wantedBalancerCalls, spyBalancer.calls)
		})
	}
}
