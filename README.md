# go-parallel

Package `parallel` provides a construct to simplify parallel for loop execution in Go,
inspired by OpenMP's "parallel for" pragmas in C, C++ and Fortran.

```go
/// without parallel construct ~290ms
for i := 0; i < N; i++ {
    outputArray[i] = sinc(inputArray[i] * math.Pi)
}

// with parallel construct ~90ms on 4 cores
parallel.For(N, func(i, _ int) {
    outputArray[i] = sinc(inputArray[i] * math.Pi)
})
```

For more, visit the [GoDoc](https://godoc.org/github.com/dgravesa/go-parallel/parallel)

## Installation

```
go get -v github.com/dgravesa/go-parallel
```

## Examples

* [For loop](https://godoc.org/github.com/dgravesa/go-parallel/parallel#For)
* [For loop including goroutine ID](https://godoc.org/github.com/dgravesa/go-parallel/parallel#ForWithGrID)
* [go-modalclust](https://github.com/dgravesa/go-modalclust/blob/master/pkg/modalclust/mac.go#L30),
a Go-based implementation of a modal clustering algorithm.
