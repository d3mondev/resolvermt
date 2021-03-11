package multidns

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	mockRecords := []Record{
		{"test", TypeA, "127.0.0.1"},
	}

	mockA := &dns.A{}
	mockA.Hdr.Rrtype = dns.TypeA
	mockA.A = net.ParseIP("127.0.0.1")

	mockMsg := &dns.Msg{}
	mockMsg.Answer = []dns.RR{mockA}

	mockMsgErr := &dns.Msg{}
	mockMsgErr.Rcode = dns.RcodeServerFailure

	testTable := []struct {
		name    string
		retries int
		errors  int
		msg     *dns.Msg
		want    []Record
	}{
		{name: "Single Answer", retries: 3, errors: 0, msg: mockMsg, want: mockRecords},
		{name: "Retry", retries: 3, errors: 2, msg: mockMsg, want: mockRecords},
		{name: "Max Retry", retries: 1, errors: 1, msg: mockMsg, want: []Record{}},
		{name: "DNS Error", retries: 3, errors: 0, msg: mockMsgErr, want: []Record{}},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRequest := NewMockrequest(ctrl)
			mockServer := newRateLimitedServer("8.8.8.8", 10)
			mockBalancer := NewMockbalancer(ctrl)
			mockParser := NewMockparser(ctrl)
			mockNewSender := func(query string, rrtype RRtype) sender {
				return mockRequest
			}

			// Send should return errors up until the number of retries has been reached
			for i := 0; i < test.errors && i < test.retries; i++ {
				mockBalancer.EXPECT().Next().Return(mockServer)
				mockRequest.EXPECT().Send(gomock.Any()).Return(test.msg, time.Duration(0), errors.New("error"))
			}

			// DNS request successful
			if test.retries > test.errors {
				mockBalancer.EXPECT().Next().Return(mockServer)
				mockRequest.EXPECT().Send(gomock.Any()).Return(test.msg, time.Duration(0), nil)

				if test.msg != mockMsgErr {
					mockParser.EXPECT().Parse("test", mockMsg).Return([]Record{{Question: "test", Type: TypeA, Answer: "127.0.0.1"}})
				}
			}

			resolver := newResolver(test.retries, mockNewSender, mockBalancer, mockParser)
			channel := make(chan []Record, 10)
			resolver.Resolve("test", TypeA, channel)
			got := <-channel

			assert.EqualValues(t, test.want, got, test.name)
		})
	}
}
