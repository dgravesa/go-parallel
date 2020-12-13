package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/dgravesa/go-parallel/parallel"
)

func main() {
	var N int
	var runSerial bool

	flag.IntVar(&N, "N", 1000000, "number of items in input array")
	flag.BoolVar(&runSerial, "Serial", false, "execute loop serially")
	flag.Parse()

	// create input array
	x := make([]int, N)
	for i := 0; i < N; i++ {
		x[i] = i
	}

	expectedSum := N * (N - 1) / 2
	fmt.Println("expected sum:", expectedSum)

	t1 := time.Now()

	sum := 0
	if runSerial {
		for i := 0; i < N; i++ {
			sum += x[i]
		}
	} else {
		// execute partial sums
		psums := make([]int, parallel.DefaultNumGoroutines())

		parallel.For(N, func(i, grID int) {
			psums[grID] += x[i]
		})

		fmt.Println("partial sums:", psums)

		for _, psum := range psums {
			sum += psum
		}
	}

	t2 := time.Now()

	fmt.Println("result sum:", sum)

	fmt.Println("execution time:", t2.Sub(t1))
}
