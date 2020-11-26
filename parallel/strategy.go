package parallel

import (
	"math"
	"runtime"
	"sync"
)

// Strategy contains the parallel execution parameters.
type Strategy struct {
	numGoroutines int
	lockOSThreads bool
}

// NewStrategy returns a new parallel strategy
func NewStrategy() *Strategy {
	s := new(Strategy)
	s.numGoroutines = 1
	s.lockOSThreads = false
	return s
}

// DefaultStrategy returns the default parallel strategy
func DefaultStrategy() *Strategy {
	s := new(Strategy)
	s.numGoroutines = runtime.GOMAXPROCS(0)
	s.lockOSThreads = true
	return s
}

// NumGoroutines returns the number of goroutines that a strategy will use
func (s *Strategy) NumGoroutines() int {
	return s.numGoroutines
}

// WithNumGoroutines sets the number of goroutines for a parallel strategy
func (s *Strategy) WithNumGoroutines(numGoroutines int) *Strategy {
	s.numGoroutines = numGoroutines
	return s
}

// WithCPUProportion sets the number of goroutines based on a proportion of number of CPUs
func (s *Strategy) WithCPUProportion(p float64) *Strategy {
	numCPU := runtime.NumCPU()
	pCPU := p * float64(numCPU)
	s.numGoroutines = int(math.Max(pCPU, 1.0))
	return s
}

// WithOSThreadLock sets whether or not to lock loop goroutines to OS threads
func (s *Strategy) WithOSThreadLock(lockOSThreads bool) *Strategy {
	s.lockOSThreads = lockOSThreads
	return s
}

// For executes a loop in parallel from i = 0 while i < N using the given strategy
func (s *Strategy) For(N int, loopBody func(i int)) {
	loopBodyWithGrID := func(i, _ int) {
		loopBody(i)
	}

	s.ForWithGrID(N, loopBodyWithGrID)
}

// ForWithGrID executes a loop in parallel from i = 0 while i < N using the given strategy
func (s *Strategy) ForWithGrID(N int, loopBody func(i, grID int)) {
	var wg sync.WaitGroup
	wg.Add(s.numGoroutines)

	// launch goroutines
	for grID := 0; grID < s.numGoroutines; grID++ {
		go func(grID int) {
			defer wg.Done()

			if s.lockOSThreads {
				runtime.LockOSThread()
			}

			// execute goroutine's index block
			first, last := s.grIndexBlock(grID, N)
			for i := first; i < last; i++ {
				loopBody(i, grID)
			}

			if s.lockOSThreads {
				runtime.UnlockOSThread()
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
