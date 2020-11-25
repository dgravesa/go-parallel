package parallel

import "testing"

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
	For(len(slice), func(i int) {
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
	numGR := DefaultNumGoroutines()
	partialSums := make([]int, numGR)

	// act
	ForWithGrID(N, func(i, grID int) {
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
	s := WithNumGoroutines(numGoroutines)

	// assert
	if s.numGoroutines != numGoroutines {
		t.Errorf("expected %d, actual %d\n", numGoroutines, s.numGoroutines)
	}
}
