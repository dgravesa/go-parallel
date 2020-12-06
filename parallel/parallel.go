package parallel

// For executes N iterations of a function body divided equally among GOMAXPROCS goroutines.
// This function correlates directly to a for loop of the form:
//
// 		for i := 0; i < N; i++ {
//			loopBody(i)
// 		}
//
// Replacing existing for loops with this construct may accelerate parallelizable workloads.
func For(N int, loopBody func(i int)) {
	NewExecutor().For(N, loopBody)
}

// ForWithGrID executes N iterations of a function body divided equally among GOMAXPROCS goroutines.
// Unlike For, ForWithGrID also incorporates a grID argument that may be used in the loop body.
// The grID argument is the goroutine ID and may be used for a partial reduction at the goroutine level.
// Goroutine IDs range from 0 to NumGoroutines - 1.
//
// Replacing existing for loops with this construct may accelerate parallelizable workloads.
func ForWithGrID(N int, loopBody func(i, grID int)) {
	NewExecutor().ForWithGrID(N, loopBody)
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
