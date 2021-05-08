package parallel_test

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/dgravesa/go-parallel/parallel"
)

func Example() {
	// initialize input array of N values, not included in runtime
	N := 1000
	inputArray := make([]float64, N)
	outputArray := make([]float64, N)
	for i := 0; i < N; i++ {
		inputArray[i] = 10 * (rand.Float64() - 0.5) // -5 to 5
	}

	t1 := time.Now()

	// run loop serially
	for i := 0; i < N; i++ {
		outputArray[i] = sinc(inputArray[i] * math.Pi)
	}

	t2 := time.Now()

	// run loop in parallel
	parallel.For(N, func(i, _ int) {
		outputArray[i] = sinc(inputArray[i] * math.Pi)
	})

	t3 := time.Now()

	fmt.Printf("serial: %v\n", t2.Sub(t1))
	fmt.Printf("parallel: %v (%d goroutines)\n", t3.Sub(t2), parallel.DefaultNumGoroutines())
}

func sinc(x float64) float64 {
	time.Sleep(1 * time.Millisecond) // simulate a longer running operation

	if x == 0.0 {
		return 1.0
	}
	return math.Sin(x) / x
}
