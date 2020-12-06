package parallel_test

import (
	"fmt"
	"testing"

	"github.com/dgravesa/go-parallel/parallel"
)

func Test_ExecutorFor_WithNewExecutor_ComputesCorrectResult(t *testing.T) {
	// arrange
	slice := []float64{0.0, 3.75, -1.5, -2.0, 0.5, 0.75}
	expectedResult := []float64{1.0, 4.75, -0.5, -1.0, 1.5, 1.75}

	// act
	parallel.NewExecutor().For(len(slice), func(i int) {
		slice[i] += 1.0
	})

	// assert
	assertFloat64SlicesEqual(t, expectedResult, slice, "")
}

func Test_ExecutorFor_WithVaryingNumGoroutines_ComputesCorrectResult(t *testing.T) {
	// arrange
	inputArray := []float64{0.0, 3.75, -1.5, -2.0, 0.5, 0.75, 1.0}
	expectedOutput := []float64{0.0, 7.5, -3.0, -4.0, 1.0, 1.5, 2.0}
	N := len(inputArray)

	for _, numGRs := range []int{1, 2, 3} {
		actualOutput := make([]float64, N)

		// act
		parallel.NewExecutor().WithNumGoroutines(numGRs).For(N, func(i int) {
			actualOutput[i] = 2.0 * inputArray[i]
		})

		// assert
		prefix := fmt.Sprintf("%d threads) ", numGRs)
		assertFloat64SlicesEqual(t, expectedOutput, actualOutput, prefix)
	}
}

func Test_ExecutorForWithGrID_ComputesCorrectResult(t *testing.T) {
	// arrange
	inputArray := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	expectedSum := 15 * 16 / 2
	N := len(inputArray)

	for _, numGR := range []int{1, 2, 3, 4} {
		partialSums := make([]int, numGR)
		e := parallel.NewExecutor().WithNumGoroutines(numGR)

		// act
		e.ForWithGrID(N, func(i, grID int) {
			partialSums[grID] += inputArray[i]
		})

		// assert
		actualSum := 0
		for _, partialSum := range partialSums {
			actualSum += partialSum
		}
		if expectedSum != actualSum {
			t.Errorf("%d threads) expected %d, actual %d\n", numGR, expectedSum, actualSum)
		}
	}
}

func Test_ExecutorWithCPUProportion_HasAtLeastOneGoroutine(t *testing.T) {
	// arrange
	p := 0.0
	expected := 1

	// arrange / act
	e := parallel.NewExecutor().WithCPUProportion(p)

	// assert
	actual := e.NumGoroutines()
	if expected != actual {
		t.Errorf("expected %d, actual %d\n", expected, actual)
	}
}

func Test_ExecutorNumGoroutines_ReturnsExpectedResult(t *testing.T) {
	// arrange
	expected := 3
	e := parallel.NewExecutor().WithNumGoroutines(expected)

	// act
	actual := e.NumGoroutines()

	// assert
	if expected != actual {
		t.Errorf("expected %d, actual %d\n", expected, actual)
	}
}
