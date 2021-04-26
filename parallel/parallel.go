package parallel

import (
	"context"
	"runtime"
)

var defaultNumGoroutines = runtime.GOMAXPROCS(0)

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
func For(N int, loopBody func(i, grID int)) {
	NewExecutor().For(N, loopBody)
}

func ForWithContext(ctx context.Context, N int,
	loopBody func(ctx context.Context, i, grID int)) error {

	return NewExecutor().ForWithContext(ctx, N, loopBody)
}

// WithNumGoroutines returns a default executor, but using a specific number of goroutines.
func WithNumGoroutines(n int) *Executor {
	return NewExecutor().WithNumGoroutines(n)
}

// WithCPUProportion returns a default executor, but with the number of goroutines based on a
// proportion of the number of CPUs.
func WithCPUProportion(p float64) *Executor {
	return NewExecutor().WithCPUProportion(p)
}

// WithStrategy returns a default executor, but with a particular parallel strategy for execution.
// Different parallel strategies vary on how work items are distributed among goroutines.
// The strategy types are defined as constants and follow the naming convention Strategy*.
// If an unrecognized value is specified, the default contiguous blocks strategy will be used.
func WithStrategy(strategy StrategyType) *Executor {
	return NewExecutor().WithStrategy(strategy)
}

// SetDefaultNumGoroutines sets the default number of goroutines for For() and NewExecutor(). At
// start time, the default is initialized to the result of runtime.GOMAXPROCS(0). If numGoroutines
// is less than 1, the default will be set to 1.
func SetDefaultNumGoroutines(numGoroutines int) {
	defaultNumGoroutines = maxInt(numGoroutines, 1)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// DefaultNumGoroutines returns the default number of goroutines for For() and NewExecutor().
// Unless specified otherwise using SetDefaultNumGoroutines(), the default number of goroutines is
// equal to the result of runtime.GOMAXPROCS(0) at start time.
func DefaultNumGoroutines() int {
	return defaultNumGoroutines
}
