package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/dgravesa/go-parallel/parallel"
)

func main() {
	var N int
	flag.IntVar(&N, "N", 1000, "number of vector elements")
	flag.Parse()

	serialX := make([]int, N)
	parallelX := make([]int, N)

	// run computation serially
	si := time.Now()
	for i := 0; i < N; i++ {
		serialX[i] = i
	}
	sf := time.Now()
	fmt.Println("serial time:", sf.Sub(si))

	// run computation in parallel
	pi := time.Now()
	parallel.For(N, func(i int) {
		parallelX[i] = i
	})
	pf := time.Now()
	fmt.Println("parallel time:", pf.Sub(pi))
}
