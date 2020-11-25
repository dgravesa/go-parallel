package parallel

import (
	"fmt"
	"testing"
)

func Test_StrategyFor_WithNewStrategy_ComputesCorrectResult(t *testing.T) {
	// arrange
	slice := []float64{0.0, 3.75, -1.5, -2.0, 0.5, 0.75}
	expectedResult := []float64{1.0, 4.75, -0.5, -1.0, 1.5, 1.75}

	// act
	NewStrategy().For(len(slice), func(i int) {
		slice[i] += 1.0
	})

	// assert
	assertFloat64SlicesEqual(t, expectedResult, slice, "")
}

func Test_StrategyFor_WithDefaultStrategy_ComputesCorrectResult(t *testing.T) {
	// arrange
	slice := []float64{0.0, 3.75, -1.5, -2.0, 0.5, 0.75, 1.0, 4.5, -15.0}
	expectedResult := []float64{1.0, 4.75, -0.5, -1.0, 1.5, 1.75, 2.0, 5.5, -14.0}

	// act
	DefaultStrategy().For(len(slice), func(i int) {
		slice[i] += 1.0
	})

	// assert
	assertFloat64SlicesEqual(t, expectedResult, slice, "")
}

func Test_StrategyFor_WithVaryingNumGoroutines_ComputesCorrectResult(t *testing.T) {
	// arrange
	inputArray := []float64{0.0, 3.75, -1.5, -2.0, 0.5, 0.75, 1.0}
	expectedOutput := []float64{0.0, 7.5, -3.0, -4.0, 1.0, 1.5, 2.0}
	N := len(inputArray)

	for _, numGRs := range []int{1, 2, 3} {
		actualOutput := make([]float64, N)

		// act
		DefaultStrategy().WithNumGoroutines(numGRs).For(N, func(i int) {
			actualOutput[i] = 2.0 * inputArray[i]
		})

		// assert
		prefix := fmt.Sprintf("%d threads) ", numGRs)
		assertFloat64SlicesEqual(t, expectedOutput, actualOutput, prefix)
	}
}

func Test_StrategyGRIndexBlock_ReturnsCorrectRange(t *testing.T) {
	type TestCase struct {
		numGR, grID, N      int
		expStart, expFinish int
	}

	// arrange
	testCases := []TestCase{
		{1, 0, 13, 0, 13},
		{2, 0, 15, 0, 8},
		{2, 1, 15, 8, 15},
		{3, 0, 20, 0, 7},
		{3, 1, 20, 7, 14},
		{3, 2, 20, 14, 20},
		{3, 0, 1000, 0, 334},
		{3, 1, 1000, 334, 667},
		{3, 2, 1000, 667, 1000},
	}

	for i, tc := range testCases {
		// act
		s := NewStrategy().WithNumGoroutines(tc.numGR)
		actualStart, actualFinish := s.grIndexBlock(tc.grID, tc.N)

		// assert
		if tc.expStart != actualStart || tc.expFinish != actualFinish {
			t.Errorf("%d) %v expected = [%d, %d) actual [%d, %d)\n", i, tc,
				tc.expStart, tc.expFinish, actualStart, actualFinish)
		}
	}
}
