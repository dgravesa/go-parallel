// Package parallel provides a looping construct that enables parallel execution,
// inspired by OpenMP's "parallel for" pragmas in C, C++ and Fortran.
//
// Introduction
//
// The parallel package makes it easy to accelerate parallelizable loops on multicore systems:
//
// 		// loop using the parallel construct in this package,
//		// a general example that computes N outputs from a slice of N independent inputs
//
// 		parallel.For(N, func(i, _ int) {
// 			outputs[i] = computeResult(inputs[i])
// 		})
//
// At its core, the parallel package is a wrapper around a common pattern
// for parallelization in Go using sync.WaitGroup, similar to:
//
// 		// same result as before, but using sync.WaitGroup
//
// 		var wg sync.WaitGroup
// 		wg.Add(N)
// 		for i := 0; i < N; i++ {
// 			go func(i int) {
//				defer wg.Done()
// 				outputs[i] = computeResult(inputs[i])
// 			}(i)
// 		}
// 		wg.Wait()
//
// Whereas this snippet using sync.WaitGroup creates a goroutine for each loop iteration, the
// parallel construct in this package abstracts the goroutine logic and distributes the work
// automatically and intelligently among a smaller number of goroutines, minimizing the overhead
// that results from excessive goroutine lifecycling and scheduling.
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
// Use Cases
//
// In general, the parallel package may be used to distribute N loop iterations across a fixed
// number of goroutines. Some use cases include:
//
// • Accelerating embarrassingly parallel for loops, such as large vector operations.
//
// 		N := 10000000
//
// 		// for i := 0; i < N; i++ {
// 		// 	z[i] = x[i] + y[i]
// 		// }
//
//		// equivalent using parallel package
// 		parallel.For(N, func(i, _ int) {
// 			z[i] = x[i] + y[i]
// 		})
//
// • A replacement for the common sync.WaitGroup pattern:
//
// 		// var wg sync.WaitGroup
// 		// wg.Add(N)
// 		// for i := 0; i < N; i++ {
// 		// 	go func(i int) {
// 		// 		defer wg.Done()
// 		// 		executeTask(tasks[i])
// 		// 	}(i)
// 		// }
// 		// wg.Wait()
//
// 		// equivalent using parallel package
// 		parallel.WithNumGoroutines(N).For(N, func(i, _ int) {
// 			executeTask(tasks[i])
// 		})
//
// • Batching API requests. For example, if I need to make 200 independent API requests, but want
// to limit to 30 active requests at a time, I can accomplish this using the parallel package as
// follows:
//
// 		concurrency := 30
// 		numRequests := 200
//
// 		// NOTE: use StrategyAtomicCounter since API requests tend to vary in response times
// 		requestsExecutor := parallel.NewExecutor().
// 			WithStrategy(parallel.StrategyAtomicCounter).
// 			WithNumGoroutines(concurrency)
//
// 		requestsExecutor.For(numRequests, func(i, _ int) {
// 			responses[i] = executeAPIRequest(requests[i])
// 		})
//
// Best Practices - Verify Speedup
//
// • Not every loop is going to be faster simply by using the parallel execution provided by this
// package. Loops with small amounts of work
// per index (such as basic arithmetic operations) will often see no benefit from using this
// package, and may actually run slower.
//
// • A good rule of thumb is if the loop body makes at least one call to a non-inlineable
// function, it may benefit from this parallel package.
//
// • Always verify results when parallelizing loops, both for speedup and correctness.
//
// Best Practices - Selecting a Strategy
//
// • Generally, the default StrategyContiguousBlocks is recommended in all loops when each loop
// index takes less than one microsecond.
//
// • In loops where the amount of work or time varies by index, using StrategyAtomicCounter may
// help to balance work more evenly among goroutines, resulting in faster loop execution.
//
package parallel
