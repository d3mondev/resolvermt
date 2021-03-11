package multidns

import (
	"sync/atomic"
)

type roundRobinBalancer struct {
	values   []server
	curIndex int64
	count    int
}

func newRoundRobinBalancer(servers []server) *roundRobinBalancer {
	count := len(servers)

	list := roundRobinBalancer{
		values:   make([]server, count),
		curIndex: -1,
		count:    count,
	}

	copy(list.values, servers)

	return &list
}

func (s *roundRobinBalancer) Next() server {
	if len(s.values) == 0 {
		return nil
	}

	nextIndex := int(atomic.AddInt64(&s.curIndex, 1))
	server := s.values[nextIndex%s.count]

	return server
}
