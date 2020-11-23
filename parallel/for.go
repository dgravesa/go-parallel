package parallel

// For executes a loop in parallel from i = 0 while i < N
func For(N int, loopBody func(i int)) {
	Default().For(N, loopBody)
}
