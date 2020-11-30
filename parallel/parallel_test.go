package parallel_test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/dgravesa/go-parallel/parallel"
)

func ExampleFor() {
	x := []int{1, 2, 3, 4, 5, 6, 7, 8}
	y := []int{0, 1, 0, 1, 0, 1, 0, 1}
	N := len(x)
	z := make([]int, N)

	// compute z = x * y
	parallel.For(N, func(i int) {
		z[i] = x[i] * y[i]
	})

	fmt.Println(z)
	// Output: [0 2 0 4 0 6 0 8]
}

func ExampleForWithGrID() {
	x := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	N := len(x)
	psums := make([]int, parallel.DefaultNumGoroutines())

	// compute partial sums
	parallel.ForWithGrID(N, func(i, grID int) {
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

	// compute using 3 goroutines
	parallel.WithNumGoroutines(3).For(N, func(i int) {
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

	// compute z = x * y using 70% of CPUs, minimum 1
	parallel.WithCPUProportion(0.7).For(N, func(i int) {
		floor[i] = int(x[i])
	})

	fmt.Println(floor)
	// Output: [1 2 1 5 3 9 6 6]
}

func assertFloat64SlicesEqual(t *testing.T, expected, actual []float64, prefix string) {
	if len(expected) != len(actual) {
		t.Errorf("%sslices do not have same dimension: len(expected) = %d, len(actual) = %d\n",
			prefix, len(expected), len(actual))
		return
	}

	failed := false
	for i := 0; i < len(expected); i++ {
		if expected[i] != actual[i] {
			failed = true
		}
	}
	if failed {
		t.Errorf("%sslices do not match: expected %v, actual %v\n", prefix, expected, actual)
	}
}

func Test_For_ComputesCorrectResult(t *testing.T) {
	// arrange
	slice := []float64{0.0, 3.75, -1.5, -2.0, 0.5, 0.75}
	expectedResult := []float64{1.0, 4.75, -0.5, -1.0, 1.5, 1.75}

	// act
	parallel.For(len(slice), func(i int) {
		slice[i] += 1.0
	})

	// assert
	assertFloat64SlicesEqual(t, expectedResult, slice, "")
}

func Test_ForWithGrID_ComputesCorrectResult(t *testing.T) {
	// arrange
	inputArray := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	expectedSum := 9 * 10 / 2
	N := len(inputArray)
	numGR := parallel.DefaultNumGoroutines()
	partialSums := make([]int, numGR)

	// act
	parallel.ForWithGrID(N, func(i, grID int) {
		partialSums[grID] += inputArray[i]
	})

	// assert
	actualSum := 0
	for _, partialSum := range partialSums {
		actualSum += partialSum
	}
	if expectedSum != actualSum {
		t.Errorf("expected %d, actual %d\n", expectedSum, actualSum)
	}
}

func Test_WithNumGoroutines_ReturnsValidStrategy(t *testing.T) {
	// arrange
	numGoroutines := 3

	// act
	s := parallel.WithNumGoroutines(numGoroutines)

	// assert
	if s.NumGoroutines() != numGoroutines {
		t.Errorf("expected %d, actual %d\n", numGoroutines, s.NumGoroutines())
	}
}

func BenchmarkForSinc(b *testing.B) {
	N := 1000000
	inputArray := make([]float64, N)
	outputArray := make([]float64, N)
	for i := 0; i < N; i++ {
		inputArray[i] = 10 * (rand.Float64() - 0.5)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		parallel.For(N, func(i int) {
			xPi := inputArray[i] * math.Pi
			outputArray[i] = math.Sin(xPi) / xPi
		})
	}
}
