package parallel

import (
	"testing"
)

func Test_ExecutorGRIndexBlock_ReturnsCorrectRange(t *testing.T) {
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
		e := NewExecutor().WithNumGoroutines(tc.numGR)
		actualStart, actualFinish := e.grIndexBlock(tc.grID, tc.N)

		// assert
		if tc.expStart != actualStart || tc.expFinish != actualFinish {
			t.Errorf("%d) %v expected = [%d, %d) actual [%d, %d)\n", i, tc,
				tc.expStart, tc.expFinish, actualStart, actualFinish)
		}
	}
}
