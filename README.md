# go-parallel

Package `parallel` provides a construct to simplify parallel for loop execution in Go,
inspired by OpenMP's "parallel for" pragmas in C, C++ and Fortran.

```go
// without parallel construct ~290ms
for i := 0; i < N; i++ {
    outputArray[i] = sinc(inputArray[i] * math.Pi)
}

// with parallel construct ~90ms on 4 cores
parallel.For(N, func(i, _ int) {
    outputArray[i] = sinc(inputArray[i] * math.Pi)
})
```

## Installation

```
go get -v github.com/dgravesa/go-parallel
```

## API Reference / Examples

Visit the [GoDoc](https://godoc.org/github.com/dgravesa/go-parallel/parallel) for API reference and examples.

## About

The parallel package makes it easy to accelerate parallelizable loops on multicore systems:

```go
// loop using the parallel construct in this package,
// a general example that computes N outputs from a slice of N independent inputs

parallel.For(N, func(i, _ int) {
    outputs[i] = computeResult(inputs[i])
})
```

### Motivation

Go is designed for extreme concurrency.
Goroutines work well for this as they are considerably lightweight.
However, goroutines are not free.
On for loops with many iterations, launching one goroutine per iteration may prove to be
overkill, resulting in slower execution than if the loop were to be run serially.

```go
N := 10000000
inputArray := make([]float64, N)
outputArray := make([]float64, N)
for i := 0; i < N; i++ {
    inputArray[i] = 10 * (rand.Float64() - 0.5) // -5 to 5
}

sinc := func(x float64) float64 {
    if x == 0.0 {
        return 1.0
    }
    return math.Sin(x) / x
}

// serial
// ~290ms
for i := 0; i < N; i++ {
    outputArray[i] = sinc(inputArray[i] * math.Pi)
}

// one goroutine per iteration with 4 CPUs
// ~1.9s
for i := 0; i < N; i++ {
    go func(i int) {
        outputArray[i] = sinc(inputArray[i] * math.Pi)
    }(i)
}

// parallel package construct with 4 CPUs
// ~90ms
parallel.For(N, func(i, _ int) {
    outputArray[i] = sinc(inputArray[i] * math.Pi)
})
```

The constructs provided by the parallel package automatically handle goroutine management and
distribution of work.
By default, the number of goroutines is set to the number of CPUs.
This way, the parallel constructs avoid the overhead that results from excessive goroutine
creation and scheduling.

### Use Cases

In general, the parallel package may be used to distribute N loop iterations across a fixed
number of goroutines. Some use cases include:

* Accelerating embarrassingly parallel for loops, such as large vector operations.

```go
N := 10000000

// for i := 0; i < N; i++ {
// 	z[i] = x[i] + y[i]
// }

// equivalent using parallel package
parallel.For(N, func(i, _ int) {
    z[i] = x[i] + y[i]
})
```

* A replacement for the common sync.WaitGroup pattern:

```go
// var wg sync.WaitGroup
// wg.Add(N)
// for i := 0; i < N; i++ {
// 	go func(i int) {
// 		defer wg.Done()
// 		executeTask(tasks[i])
// 	}(i)
// }
// wg.Wait()

// equivalent using parallel package
parallel.WithNumGoroutines(N).For(N, func(i, _ int) {
    executeTask(tasks[i])
})
```

* Batching API requests. For example, if I need to make 200 independent API requests, but want
to limit to 30 active requests at a time, I can accomplish this using the parallel package as
follows:

```go
concurrency := 30
numRequests := 200

// NOTE: use StrategyFetchNextIndex since API requests tend to vary in response times
requestsExecutor := parallel.NewExecutor().
    WithStrategy(parallel.StrategyFetchNextIndex).
    WithNumGoroutines(concurrency)

requestsExecutor.For(numRequests, func(i, _ int) {
    responses[i] = executeAPIRequest(requests[i])
})
```

## Application Tuning

Ultimately, the optimal strategy and number of goroutines will vary from loop to loop.
This section provides a few rules of thumb.
Although not always necessary, generating a [CPU profile](https://pkg.go.dev/runtime/pprof) may provide additional information for better application performance.

It's important to verify speedup over the existing serial implementation.
In many cases, loops will not be large enough to benefit from the parallel execution provided by this package.

### Selecting a Strategy

| Strategy | Use Cases |
| -------- | --------- |
| StrategyPreassignIndices | Each loop iteration takes less than one microsecond. |
| StrategyFetchNextIndex | Some or all loop iterations take longer than one microsecond. |

### Selecting number of goroutines

* For compute-bound loops, the optimal number of goroutines is typically equal to or slightly less than the number of CPUs.
By default, parallel execution will execute a number of goroutines equal to the number of CPUs. To use slightly less than this,
specify either `WithCPUProportion(p)` with *p < 1.0* or `WithNumGoroutines(n)`.
* For network-bound loops, the optimal number of goroutines may depend on the network bandwidth required for each iteration, but will often be more than the number of CPUs. In this case, `WithNumGoroutines(n)` should be tested with increasing values for *n* until an optimal value is found.
