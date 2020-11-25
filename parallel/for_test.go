package parallel

import "testing"

func Test_For_ComputesCorrectResult(t *testing.T) {
	// arrange
	slice := []float64{0.0, 3.75, -1.5, -2.0, 0.5, 0.75}
	expectedResult := []float64{1.0, 4.75, -0.5, -1.0, 1.5, 1.75}

	// act
	For(len(slice), func(i int) {
		slice[i] += 1.0
	})

	// assert
	for i := 0; i < len(slice); i++ {
		failed := false
		if slice[i] != expectedResult[i] {
			failed = true
		}

		if failed {
			t.Errorf("expected %v, actual %v\n", expectedResult, slice)
		}
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
