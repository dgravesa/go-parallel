package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/dgravesa/go-parallel/parallel"
)

func main() {
	var N int
	var runSerial bool

	flag.IntVar(&N, "N", 1000000, "number of vector elements")
	flag.BoolVar(&runSerial, "Serial", false, "run non-parallelized for loop instead")
	flag.Parse()

	x := make([]int, N)

	t1 := time.Now()
	if runSerial {
		// run computation serially
		for i := 0; i < N; i++ {
			x[i] = i
		}
	} else {
		// run computation in parallel
		parallel.For(context.TODO(), N, func(pctx *parallel.Context) {
			i := pctx.Index()
			x[i] = i
		})
	}
	t2 := time.Now()
	fmt.Println(t2.Sub(t1))
}
