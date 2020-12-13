package parallel_test

import (
	"fmt"

	"github.com/dgravesa/go-parallel/parallel"
)

func ExampleFor_basic() {
	x := []int{1, 2, 3, 4, 5, 6, 7, 8}
	y := []int{0, 1, 0, 1, 0, 1, 0, 1}
	N := len(x)
	z := make([]int, N)

	// compute z = x * y
	parallel.For(N, func(i, _ int) {
		z[i] = x[i] * y[i]
	})

	fmt.Println(z)
	// Output: [0 2 0 4 0 6 0 8]
}

func ExampleFor_goroutineID() {
	x := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	N := len(x)
	psums := make([]int, parallel.DefaultNumGoroutines())

	// compute partial sums
	parallel.For(N, func(i, grID int) {
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

func ExampleExecutor_WithNumGoroutines() {
	x := []int{1, 2, 3, 4, 5, 6, 7}
	N := len(x)
	isEven := make([]bool, N)

	// compute using 3 goroutines
	parallel.WithNumGoroutines(3).For(N, func(i, _ int) {
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

func ExampleExecutor_WithCPUProportion() {
	x := []float64{1.2, 2.0, 1.9, 5.5, 3.4, 9.3, 6.4, 6.6}
	N := len(x)
	floor := make([]int, N)

	// compute z = x * y using 70% of CPUs, minimum 1
	parallel.WithCPUProportion(0.7).For(N, func(i, _ int) {
		floor[i] = int(x[i])
	})

	fmt.Println(floor)
	// Output: [1 2 1 5 3 9 6 6]
}