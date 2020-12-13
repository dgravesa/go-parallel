package parallel

// For executes N iterations of a function body divided equally among a number of goroutines.
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

// WithNumGoroutines returns a default executor, but using a specific number of goroutines.
func WithNumGoroutines(n int) *Executor {
	return NewExecutor().WithNumGoroutines(n)
}

// WithCPUProportion returns a default executor, but with the number of goroutines based on a
// proportion of the number of CPUs.
func WithCPUProportion(p float64) *Executor {
	return NewExecutor().WithCPUProportion(p)
}

// DefaultNumGoroutines returns the default number of goroutines.
func DefaultNumGoroutines() int {
	return NewExecutor().NumGoroutines()
}
