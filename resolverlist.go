package multidns

import "sync/atomic"

type resolverList interface {
	GetNextResolver() resolver
	Count() int
}

type resolverListRoundRobin struct {
	resolvers     []resolver
	resolverIndex int32
	count         int
}

func newResolverListRoundRobin(resolvers []string, queriesPerSecond int) resolverList {
	if resolvers == nil || len(resolvers) == 0 {
		return nil
	}

	resolverList := resolverListRoundRobin{}

	for _, server := range resolvers {
		r := newResolver(server, queriesPerSecond)
		resolverList.resolvers = append(resolverList.resolvers, r)
	}

	resolverList.resolverIndex = -1
	resolverList.count = len(resolverList.resolvers)

	return &resolverList
}

func (s *resolverListRoundRobin) Count() int {
	return s.count
}

func (s *resolverListRoundRobin) GetNextResolver() resolver {
	nextIndex := int(atomic.AddInt32(&s.resolverIndex, 1))
	resolver := s.resolvers[nextIndex%s.count]

	return resolver
}
