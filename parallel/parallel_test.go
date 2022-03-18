package parallel_test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/dgravesa/go-parallel/parallel"
)

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

func Test_For_Basic_ComputesCorrectResult(t *testing.T) {
	// arrange
	slice := []float64{0.0, 3.75, -1.5, -2.0, 0.5, 0.75}
	expectedResult := []float64{1.0, 4.75, -0.5, -1.0, 1.5, 1.75}

	// act
	parallel.For(len(slice), func(i, _ int) {
		slice[i] += 1.0
	})

	// assert
	assertFloat64SlicesEqual(t, expectedResult, slice, "")
}

func Test_For_WithGrID_ComputesCorrectResult(t *testing.T) {
	// arrange
	inputArray := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	expectedSum := 9 * 10 / 2
	N := len(inputArray)
	numGR := parallel.DefaultNumGoroutines()
	partialSums := make([]int, numGR)

	// act
	parallel.For(N, func(i, grID int) {
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

func Test_WithNumGoroutines_ReturnsValidExecutor(t *testing.T) {
	// arrange
	numGoroutines := 3

	// act
	e := parallel.WithNumGoroutines(numGoroutines)

	// assert
	if e.NumGoroutines() != numGoroutines {
		t.Errorf("expected %d, actual %d\n", numGoroutines, e.NumGoroutines())
	}
}

func Test_WithStrategy_ReturnsValidExecutor(t *testing.T) {
	// arrange
	strategy := parallel.StrategyFetchNextIndex
	resultArray := []float64{0, 0, 0, 0, 0, 0, 0}
	expectedResult := []float64{0.5, 1.5, 2.5, 3.5, 4.5, 5.5, 6.5}
	N := len(resultArray)

	// act
	parallel.WithStrategy(strategy).For(N, func(i, _ int) {
		resultArray[i] = float64(i) + 0.5
	})

	// assert
	assertFloat64SlicesEqual(t, expectedResult, resultArray, "")
}

func Test_SetDefaultNumGoroutines_WithPositiveInteger_SetsToThatInteger(t *testing.T) {
	// arrange
	expectedDefaultNumGR := 27

	// act
	parallel.SetDefaultNumGoroutines(expectedDefaultNumGR)

	// assert
	actualDefaultNumGR := parallel.DefaultNumGoroutines()
	if expectedDefaultNumGR != actualDefaultNumGR {
		t.Errorf("expected %d, actual %d\n", expectedDefaultNumGR, actualDefaultNumGR)
	}
}

func Test_SetDefaultNumGoroutines_WithNegativeInteger_SetsToOne(t *testing.T) {
	// arrange
	defaultNumGRArg := -12
	expectedDefaultNumGR := 1

	// act
	parallel.SetDefaultNumGoroutines(defaultNumGRArg)

	// assert
	actualDefaultNumGR := parallel.DefaultNumGoroutines()
	if expectedDefaultNumGR != actualDefaultNumGR {
		t.Errorf("expected %d, actual %d\n", expectedDefaultNumGR, actualDefaultNumGR)
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
		parallel.For(N, func(i, _ int) {
			xPi := inputArray[i] * math.Pi
			outputArray[i] = math.Sin(xPi) / xPi
		})
	}
}
