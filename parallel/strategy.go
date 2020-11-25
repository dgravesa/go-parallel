package parallel

import (
	"runtime"
	"sync"
)

// Strategy contains the parallel execution parameters.
type Strategy struct {
	numGoroutines int
}

// NewStrategy returns a new parallel strategy
func NewStrategy() *Strategy {
	s := new(Strategy)
	s.numGoroutines = 1
	return s
}

// DefaultStrategy returns the default parallel strategy
func DefaultStrategy() *Strategy {
	s := new(Strategy)
	s.numGoroutines = runtime.GOMAXPROCS(0)
	return s
}

// WithNumGoroutines sets the number of goroutines for a parallel strategy
func (s *Strategy) WithNumGoroutines(numGoroutines int) *Strategy {
	s.numGoroutines = numGoroutines
	return s
}

// For executes a loop in parallel from i = 0 while i < N using the given strategy
func (s *Strategy) For(N int, loopBody func(i int)) {
	var wg sync.WaitGroup
	wg.Add(s.numGoroutines)

	// launch goroutines
	for grID := 0; grID < s.numGoroutines; grID++ {
		go func(grID int) {
			defer wg.Done()
			first, last := s.grIndexBlock(grID, N)
			for i := first; i < last; i++ {
				loopBody(i)
			}
		}(grID)
	}

	wg.Wait()
}

// grIndexBlock computes the contiguous index range for a goroutine with given ID
func (s *Strategy) grIndexBlock(grID, N int) (int, int) {
	div := N / s.numGoroutines
	mod := N % s.numGoroutines

	numWorkItems := div
	if grID < mod {
		numWorkItems++
	}

	firstIndex := grID*div + minInt(grID, mod)
	lastIndex := firstIndex + numWorkItems

	return firstIndex, lastIndex
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
