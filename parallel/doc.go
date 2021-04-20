// Package parallel provides looping constructs that enable parallel execution,
// inspired by OpenMP.
//
// Introduction
//
// At its core, the parallel package is a wrapper around a common pattern
// for parallelization in Go using sync.WaitGroup:
//
// 		N := len(inputs)
// 		outputs := make([]result, N)
//
// 		// using sync.WaitGroup
// 		var wg sync.WaitGroup
// 		wg.Add(N)
// 		for i := 0; i < N; i++ {
// 			go func(i int) {
//				// some arbitrary time-consuming task
// 				outputs[i] = processInput(inputs[i])
// 			}(i)
// 		}
// 		wg.Wait()
//
// 		// same loop, but using the parallel package
// 		parallel.For(N, func(i, _ int) {
//			// some arbitrary time-consuming task
// 			outputs[i] = processInput(inputs[i])
// 		})
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
// By default, the number of goroutines is based on the number of CPUs.
// This way, the parallel constructs avoid the overhead that results from excessive goroutine
// creation and scheduling.
//
// Best Practices
//
// 1. For workloads where the time in each iteration tends to vary,
// use the atomic counter strategy instead of the default contiguous blocks.
//
// 2. The parallel package may be used as a convenience for managing batches of API requests.
// For example, an array of 80 API requests could be limited to 10 concurrent requests as follows:
//
// 		N := len(requests) // N = 80
// 		responses := make([]Response, N)
//
// 		// NOTE: atomic counter strategy since API response times tend to vary.
// 		executor := parallel.WithNumGoroutines(10).WithStrategy(parallel.StrategyAtomicCounter)
//
// 		// execute up to 10 requests in parallel at a time
// 		executor.For(N, func(i, _ int) {
// 			responses[i] = executeRequest(requests[i])
// 		})
//
package parallel
