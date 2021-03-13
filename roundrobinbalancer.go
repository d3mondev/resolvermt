package fastdns

import (
	"sync/atomic"
)

type roundRobinBalancer struct {
	servers  []server
	curIndex int64
	count    int
}

func newRoundRobinBalancer(servers []server) *roundRobinBalancer {
	count := len(servers)

	list := roundRobinBalancer{
		servers:  make([]server, count),
		curIndex: -1,
		count:    count,
	}

	copy(list.servers, servers)

	return &list
}

func (s *roundRobinBalancer) Next() server {
	if len(s.servers) == 0 {
		return nil
	}

	nextIndex := int(atomic.AddInt64(&s.curIndex, 1))
	server := s.servers[nextIndex%s.count]

	return server
}
