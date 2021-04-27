package parallel_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgravesa/go-parallel/parallel"
)

func ExampleFor_basic() {
	x := []int{1, 2, 3, 4, 5, 6, 7, 8}
	y := []int{0, 1, 0, 1, 0, 1, 0, 1}
	N := len(x)
	z := make([]int, N)

	ctx := context.Background()

	// compute z = x * y
	parallel.For(ctx, N, func(pctx *parallel.Context) {
		i := pctx.Index()
		z[i] = x[i] * y[i]
	})

	fmt.Println(z)
	// Output: [0 2 0 4 0 6 0 8]
}

func ExampleFor_goroutineID() {
	x := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	N := len(x)
	psums := make([]int, parallel.DefaultNumGoroutines())

	ctx := context.Background()

	// compute partial sums
	parallel.For(ctx, N, func(pctx *parallel.Context) {
		i := pctx.Index()
		grID := pctx.GoroutineID()

		psums[grID] += x[i]
	})

	// compute total sum
	sum := 0
	for _, psum := range psums {
		sum += psum
	}
	fmt.Println(sum)
	// Output: 55
}

func ExampleWithNumGoroutines() {
	x := []int{1, 2, 3, 4, 5, 6, 7}
	N := len(x)
	isEven := make([]bool, N)

	ctx := context.Background()

	// compute using 3 goroutines
	parallel.WithNumGoroutines(3).For(ctx, N, func(pctx *parallel.Context) {
		i := pctx.Index()

		mod := x[i] % 2

		if mod == 1 {
			isEven[i] = false
		} else {
			isEven[i] = true
		}
	})

	fmt.Println(isEven)
	// Output: [false true false true false true false]
}

func ExampleWithCPUProportion() {
	x := []float64{1.2, 2.0, 1.9, 5.5, 3.4, 9.3, 6.4, 6.6}
	N := len(x)
	floor := make([]int, N)

	ctx := context.Background()

	// compute z = x * y using 70% of CPUs, minimum 1
	parallel.WithCPUProportion(0.7).For(ctx, N, func(pctx *parallel.Context) {
		i := pctx.Index()

		floor[i] = int(x[i])
	})

	fmt.Println(floor)
	// Output: [1 2 1 5 3 9 6 6]
}

func ExampleFor_timeout() {
	// max iteration time is 3 milliseconds
	sleepTimeMicros := []time.Duration{100, 600, 200, 100, 200, 50, 3000, 30, 10, 200, 30}
	N := len(sleepTimeMicros)

	// timeout at 1 millisecond
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Microsecond)
	defer cancel()

	err := parallel.For(ctx, N, func(pctx *parallel.Context) {
		i := pctx.Index()

		thisIterationDuration := time.Duration(sleepTimeMicros[i]) * time.Microsecond

		select {
		case <-time.After(thisIterationDuration):
			// this loop sleep iteration completed
		case <-pctx.Done():
			// deadline reached
		}
	})

	fmt.Println(errors.Is(err, context.DeadlineExceeded))
	// Output: true
}
