package parallel

// For executes a loop in parallel from i = 0 while i < N
func For(N int, loopBody func(i int)) {
	DefaultStrategy().For(N, loopBody)
}

// WithNumGoroutines returns a default strategy, but using the specifiec number of goroutines.
func WithNumGoroutines(n int) *Strategy {
	return DefaultStrategy().WithNumGoroutines(n)
}
