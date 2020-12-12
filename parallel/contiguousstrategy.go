package parallel

import "sync"

type contiguousBlocksStrategy struct{}

func (contiguousBlocksStrategy) executeFor(numGR int, N int, loopBody func(i, grID int)) {
	var wg sync.WaitGroup
	wg.Add(numGR)

	// launch goroutines
	for grID := 0; grID < numGR; grID++ {
		go func(grID int) {
			defer wg.Done()
			first, last := grIndexBlock(numGR, grID, N)
			for i := first; i < last; i++ {
				loopBody(i, grID)
			}
		}(grID)
	}

	wg.Wait()
}

// grIndexBlock computes the contiguous index range for a goroutine with given ID
func grIndexBlock(numGR, grID, N int) (int, int) {
	div := N / numGR
	mod := N % numGR

	numWorkItems := div
	if grID < mod {
		numWorkItems++
	}

	firstIndex := grID*div + minInt(grID, mod)
	lastIndex := firstIndex + numWorkItems

	return firstIndex, lastIndex
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
