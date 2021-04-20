package parallel

import (
	"math"
	"runtime"
)

// Executor is the core type used to execute parallel loops.
// New instances are created using parallel.NewExecutor().
type Executor struct {
	numGoroutines    int
	parallelStrategy strategy
}

// NewExecutor returns a new parallel executor instance.
// The default number of goroutines is equal to GOMAXPROCS.
func NewExecutor() *Executor {
	e := new(Executor)
	e.numGoroutines = runtime.GOMAXPROCS(0)
	e.parallelStrategy = defaultStrategy()
	return e
}

// NumGoroutines returns the number of goroutines that an executor is configured to use.
func (e *Executor) NumGoroutines() int {
	return e.numGoroutines
}

// WithNumGoroutines sets the number of goroutines for a parallel executor to use.
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

// WithStrategy sets the parallel strategy for execution.
// Different parallel strategies vary on how work items are distributed among goroutines.
// The strategy types are defined as constants and follow the pattern parallel.Strategy*.
// If an unrecognized value is specified, a default strategy will be chosen.
func (e *Executor) WithStrategy(strategy StrategyType) *Executor {
	switch strategy {
	case StrategyContiguousBlocks:
		e.parallelStrategy = new(contiguousBlocksStrategy)
	case StrategyAtomicCounter:
		e.parallelStrategy = new(atomicCounterStrategy)
	default:
		e.parallelStrategy = defaultStrategy()
	}
	return e
}

// For executes N iterations of a function body, where the iterations are parallelized among a
// number of goroutines.
// Replacing existing for loops with this construct may accelerate parallelizable workloads.
// The first argument to the loop body function is the loop iteration index.
// If only this argument is used, then this function correlates directly to a for loop of the form:
//
//		for i := 0; i < N; i++ {
//			loopBody(i, _)
//		}
//
// The second argument to the loop body is the ID of the goroutine executing the loop iteration.
// Goroutine IDs range from 0 to NumGoroutines - 1.
// This ID can be used as part of the parallel logic; for example, the goroutine ID may be used
// such that each goroutine computes a partial result independently, and then a final result could
// be computed more quickly from the partial results immediately after the parallel loop.
func (e *Executor) For(N int, loopBody func(i, grID int)) {
	e.parallelStrategy.executeFor(e.numGoroutines, N, loopBody)
}
