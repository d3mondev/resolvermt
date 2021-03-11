package multidns

import (
	"sync/atomic"
)

type roundRobinList struct {
	values   []server
	curIndex int64
	count    int
}

type server interface {
	Take() string
}

func newRoundRobinList(servers []server) *roundRobinList {
	count := len(servers)

	list := roundRobinList{
		values:   make([]server, count),
		curIndex: -1,
		count:    count,
	}

	copy(list.values, servers)

	return &list
}

func (s *roundRobinList) Next() server {
	if len(s.values) == 0 {
		return nil
	}

	nextIndex := int(atomic.AddInt64(&s.curIndex, 1))
	server := s.values[nextIndex%s.count]

	return server
}
