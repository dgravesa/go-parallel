// Package parallel provides a looping construct that enables parallel execution,
// inspired by OpenMP's "parallel for" pragmas in C, C++ and Fortran.
//
// Introduction
//
// The parallel package makes it easy to accelerate parallelizable loops on multicore systems:
//
// 		// loop using the parallel package
// 		parallel.For(N, func(i, _ int) {
//			// some arbitrary time-consuming task
// 			outputs[i] = processInput(inputs[i])
// 		})
//
// At its core, the parallel package is a wrapper around a common pattern
// for parallelization in Go using sync.WaitGroup:
//
// 		N := len(inputs)
// 		outputs := make([]result, N)
//
// 		// same loop, but using sync.WaitGroup
// 		var wg sync.WaitGroup
// 		wg.Add(N)
// 		for i := 0; i < N; i++ {
// 			go func(i int) {
//				defer wg.Done()
//
//				// some arbitrary time-consuming task
// 				outputs[i] = processInput(inputs[i])
// 			}(i)
// 		}
// 		wg.Wait()
//
// Motivation
//
// Go is designed for extreme concurrency.
// Goroutines work well for this as they are considerably lightweight.
// However, goroutines are not free.
// On for loops with many iterations, launching one goroutine per iteration may prove to be
// overkill, resulting in slower execution than if the loop were to be run serially.
//
//		N := 10000000
//		inputArray := make([]float64, N)
//		outputArray := make([]float64, N)
//		for i := 0; i < N; i++ {
//			inputArray[i] = 10 * (rand.Float64() - 0.5) // -5 to 5
//		}
//
// 		sinc := func(x float64) float64 {
// 			if x == 0.0 {
// 				return 1.0
// 			}
// 			return math.Sin(x) / x
// 		}
//
//		// serial
//		// ~300ms
// 		for i := 0; i < N; i++ {
// 			outputArray[i] = sinc(inputArray[i] * math.Pi)
// 		}
//
//		// one goroutine per iteration with 4 CPUs
//		// ~2.4s
//		for i := 0; i < N; i++ {
//			go func(i int) {
//				outputArray[i] = sinc(inputArray[i] * math.Pi)
//			}(i)
//		}
//
// 		// parallel package construct with 4 CPUs
//		// ~130ms
// 		parallel.For(N, func(i, _ int) {
// 			outputArray[i] = sinc(inputArray[i] * math.Pi)
// 		})
//
// The constructs provided by the parallel package automatically handle goroutine management and
// distribution of work.
// By default, the number of goroutines is set to the number of CPUs.
// This way, the parallel constructs avoid the overhead that results from excessive goroutine
// creation and scheduling.
//
// Best Practices
//
// - In loops where the amount of work or time varies by iteration,
// use the atomic counter strategy instead of the default contiguous blocks.
//
package parallel
