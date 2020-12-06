# go-parallel

Go constructs that may enable parallel execution!

```go
/// without parallel construct ~300ms
for i := 0; i < N; i++ {
    outputArray[i] = sinc(inputArray[i] * math.Pi)
}

// with parallel construct ~130ms on 4 cores
parallel.For(N, func(i int) {
    outputArray[i] = sinc(inputArray[i] * math.Pi)
})
```

For more, visit the [GoDoc](https://godoc.org/github.com/dgravesa/go-parallel/parallel)

## Basic Usage

* [For loop](https://godoc.org/github.com/dgravesa/go-parallel/parallel#For)
* [For loop including goroutine ID](https://godoc.org/github.com/dgravesa/go-parallel/parallel#ForWithGrID)