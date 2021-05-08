package parallel

import (
	"sync/atomic"
)

type atomicCounterStrategy struct {
	counter int64
}

func newAtomicCounterStrategy() Strategy {
	return &atomicCounterStrategy{
		counter: -1, // first receiver will increment atomically and receive 0
	}
}

func (s *atomicCounterStrategy) IndexGenerator(_, _, _ int) IndexGenerator {
	return &atomicIndexGenerator{
		counterAddr: &s.counter,
	}
}

type atomicIndexGenerator struct {
	counterAddr *int64
}

func (g *atomicIndexGenerator) Next() int {
	return int(atomic.AddInt64(g.counterAddr, 1))
}
