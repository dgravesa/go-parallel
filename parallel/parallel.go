package parallel

import "runtime"

// For executes a loop in parallel from i = 0 while i < N
func For(N int, loopBody func(i int)) {
	DefaultStrategy().For(N, loopBody)
}

// ForWithGrID executes a loop in parallel from i = 0 while i < N
func ForWithGrID(N int, loopBody func(i, grID int)) {
	DefaultStrategy().ForWithGrID(N, loopBody)
}

// WithNumGoroutines returns a default strategy, but using the specifiec number of goroutines.
func WithNumGoroutines(n int) *Strategy {
	return DefaultStrategy().WithNumGoroutines(n)
}

// DefaultNumGoroutines returns the default number of goroutines.
func DefaultNumGoroutines() int {
	return runtime.GOMAXPROCS(0)
}
