// Package parallel provides a looping construct that enables parallel execution,
// inspired by OpenMP's "parallel for" pragmas in C, C++ and Fortran.
//
//	// without parallel construct ~290ms
//	for i := 0; i < N; i++ {
//	    outputArray[i] = sinc(inputArray[i] * math.Pi)
//	}
//
//	// with parallel construct ~90ms on 4 cores
//	parallel.For(N, func(i, _ int) {
//		outputArray[i] = sinc(inputArray[i] * math.Pi)
//	})
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
// • Test with a varying number of goroutines. In many cases, the optimal number of goroutines may
// be less than the number of available CPU cores.
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
