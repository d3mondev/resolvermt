package fastdns

import (
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
	"go.uber.org/ratelimit"
)

type rateLimitedServer struct {
	limiter    ratelimit.Limiter
	client     *dns.Client
	ipAddrPort string

	random *rand.Rand
	mutex  sync.Mutex
}

type server interface {
	Query(query string, rrtype RRtype) (*dns.Msg, time.Duration, error)
}

func newRateLimitedServerList(ipAddrPort []string, queriesPerSecond int) []server {
	list := make([]server, len(ipAddrPort))

	for i := range ipAddrPort {
		list[i] = newRateLimitedServer(ipAddrPort[i], queriesPerSecond)
	}

	return list
}

func newRateLimitedServer(ipAddrPort string, queriesPerSecond int) *rateLimitedServer {
	if !strings.Contains(ipAddrPort, ":") {
		ipAddrPort += ":53"
	}

	server := rateLimitedServer{}
	server.limiter = ratelimit.New(queriesPerSecond, ratelimit.WithoutSlack)
	server.ipAddrPort = ipAddrPort
	server.random = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	server.client = new(dns.Client)
	server.client.Dialer = &net.Dialer{Timeout: 2000 * time.Millisecond}

	return &server
}

func (s *rateLimitedServer) Query(query string, rrtype RRtype) (*dns.Msg, time.Duration, error) {
	msg := new(dns.Msg)
	msg.Id = s.newID()
	msg.RecursionDesired = true
	msg.Question = make([]dns.Question, 1)
	msg.Question[0] = dns.Question{Name: query + ".", Qtype: uint16(rrtype), Qclass: dns.ClassINET}

	s.limiter.Take()

	return s.client.Exchange(msg, s.ipAddrPort)
}

func (s *rateLimitedServer) newID() uint16 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return uint16(s.random.Int31())
}
