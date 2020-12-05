// Package parallel provides constructs that enable parallel execution.
// Using this package may help to accelerate loop bottlenecks, often with minimal-to-no
// refactor required.
//
// Motivation
//
// Go is designed for extreme concurrency.
// Goroutines work well for this as they are considerably lightweight.
// However, goroutines are not free.
// For for loops with many iterations, launching a goroutine per iteration may prove to be
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
// The constructs provided by the parallel package automatically handle goroutine creation and
// assignment of work.
// By limiting the number of goroutines based on the number of CPUs,
// the parallel constructs avoid the overhead that results from excessive goroutine creation.
//
// 		// parallel package construct with 4 CPUs
//		// ~130ms
// 		parallel.For(N, func(i int) {
// 			outputArray[i] = sinc(inputArray[i] * math.Pi)
// 		})
package parallel
