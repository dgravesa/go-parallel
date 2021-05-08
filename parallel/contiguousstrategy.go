package parallel

type contiguousBlocksStrategy struct{}

func newContiguousBlocksStrategy() Strategy {
	return &contiguousBlocksStrategy{}
}

func (s *contiguousBlocksStrategy) IndexGenerator(numGR, grID, N int) IndexGenerator {
	startIndex, stopIndex := grIndexBlock(numGR, grID, N)

	return &contiguousIndexGenerator{
		startIndex: startIndex,
		stopIndex:  stopIndex,
		doneIndex:  N,
		nextIndex:  startIndex,
	}
}

type contiguousIndexGenerator struct {
	startIndex, stopIndex int
	doneIndex             int

	nextIndex int
}

func (g *contiguousIndexGenerator) Next() int {
	if g.nextIndex >= g.stopIndex {
		return g.doneIndex
	}

	thisIndex := g.nextIndex
	g.nextIndex++

	return thisIndex
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
