package parallel

import (
	"math"
	"runtime"
	"sync"
)

// Executor contains the parallel execution parameters.
type Executor struct {
	numGoroutines int
}

// NewExecutor returns a new parallel executor.
// The default number of goroutines is equal to GOMAXPROCS.
func NewExecutor() *Executor {
	e := new(Executor)
	e.numGoroutines = runtime.GOMAXPROCS(0)
	return e
}

// NumGoroutines returns the number of goroutines that an executor will use.
func (e *Executor) NumGoroutines() int {
	return e.numGoroutines
}

// WithNumGoroutines sets the number of goroutines for a parallel executor.
func (e *Executor) WithNumGoroutines(numGoroutines int) *Executor {
	e.numGoroutines = numGoroutines
	return e
}

// WithCPUProportion sets the number of goroutines based on a proportion of number of CPUs,
// with a minimum of 1.
func (e *Executor) WithCPUProportion(p float64) *Executor {
	numCPU := runtime.NumCPU()
	pCPU := p * float64(numCPU)
	e.numGoroutines = int(math.Max(pCPU, 1.0))
	return e
}

// For executes N iterations of a function body divided equally among a number of goroutines.
// This function correlates directly to a for loop of the form:
//
// 		for i := 0; i < N; i++ {
//			loopBody(i)
// 		}
//
// Replacing existing for loops with this construct may accelerate parallelizable workloads.
func (e *Executor) For(N int, loopBody func(i int)) {
	loopBodyWithGrID := func(i, _ int) {
		loopBody(i)
	}

	e.ForWithGrID(N, loopBodyWithGrID)
}

// ForWithGrID executes N iterations of a function body divided equally among a number of goroutines.
// Unlike For, ForWithGrID also incorporates a grID argument that may be used in the loop body.
// The grID argument is the goroutine ID and may be used for a partial reduction at the goroutine level.
// Goroutine IDs range from 0 to NumGoroutines - 1.
//
// Replacing existing for loops with this construct may accelerate parallelizable workloads.
func (e *Executor) ForWithGrID(N int, loopBody func(i, grID int)) {
	var wg sync.WaitGroup
	wg.Add(e.numGoroutines)

	// launch goroutines
	for grID := 0; grID < e.numGoroutines; grID++ {
		go func(grID int) {
			defer wg.Done()
			first, last := e.grIndexBlock(grID, N)
			for i := first; i < last; i++ {
				loopBody(i, grID)
			}
		}(grID)
	}

	wg.Wait()
}

// grIndexBlock computes the contiguous index range for a goroutine with given ID
func (e *Executor) grIndexBlock(grID, N int) (int, int) {
	div := N / e.numGoroutines
	mod := N % e.numGoroutines

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
