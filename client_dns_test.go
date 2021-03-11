package multidns

import (
	"sort"
	"sync/atomic"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockDoer struct {
	records []Record
	index   int32
}

func (s *mockDoer) Resolve(query string, rrtype RRtype) []Record {
	nextIndex := atomic.AddInt32(&s.index, 1) - 1
	record := s.records[nextIndex]

	time.Sleep(time.Duration(time.Millisecond * 10))

	return []Record{record}
}

type mockSleeper struct {
	calls int
}

func (s *mockSleeper) Sleep(t time.Duration) {
	s.calls++
}

func TestClientResolve(t *testing.T) {
	testTable := []struct {
		name       string
		concurrent int
		domains    []string
		rrtype     RRtype
		sleeper    Sleeper
		wantSleep  bool
		want       []Record
	}{
		{
			name:       "Simple",
			concurrent: 5,
			domains:    []string{"foo.bar"},
			rrtype:     TypeA,
			sleeper:    &mockSleeper{},
			wantSleep:  false,
			want: []Record{
				{
					Question: "foo.bar",
					Type:     TypeA,
					Answer:   "127.0.0.1",
				},
			},
		},
		{
			name:       "Concurrency",
			concurrent: 2,
			domains:    []string{"foo.bar", "abc.xyz"},
			rrtype:     TypeA,
			sleeper:    &mockSleeper{},
			wantSleep:  false,
			want: []Record{
				{
					Question: "foo.bar",
					Type:     TypeA,
					Answer:   "127.0.0.1",
				},
				{
					Question: "abc.xyz",
					Type:     TypeA,
					Answer:   "127.0.1.1",
				},
			},
		},
		{
			name:       "Max Concurrency",
			concurrent: 1,
			domains:    []string{"foo.bar", "abc.xyz"},
			rrtype:     TypeA,
			sleeper:    &mockSleeper{},
			wantSleep:  true,
			want: []Record{
				{
					Question: "foo.bar",
					Type:     TypeA,
					Answer:   "127.0.0.1",
				},
				{
					Question: "abc.xyz",
					Type:     TypeA,
					Answer:   "127.0.1.1",
				},
			},
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDoer := &mockDoer{records: test.want}

			client := newClientDNS(mockDoer, test.sleeper, test.concurrent)

			got := client.Resolve(test.domains, test.rrtype)

			if mockSleeper, ok := test.sleeper.(*mockSleeper); ok {
				assert.Equal(t, test.wantSleep, mockSleeper.calls > 0)
			}

			sort.SliceStable(test.want, func(i, j int) bool {
				return test.want[i].Question < test.want[j].Question
			})

			sort.SliceStable(got, func(i, j int) bool {
				return got[i].Question < got[j].Question
			})

			assert.Equal(t, test.want, got)
		})
	}
}
