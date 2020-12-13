package parallel

// StrategyType identifies a parallel strategy.
// Different parallel strategies vary on how work items are executed across goroutines.
// The strategy types are defined as constants and follow the pattern `parallel.Strategy*`.
type StrategyType int

const (
	// StrategyContiguousBlocks identifies a strategy that allocates an equal number of iterations
	// per goroutine as contiguous blocks.
	StrategyContiguousBlocks = StrategyType(iota)

	// StrategyAtomicCounter refers to a strategy that uses an atomic counter so goroutines can
	// fetch additional work items when they are ready as opposed to preallocating iterations
	// to each goroutine.
	StrategyAtomicCounter = StrategyType(iota)
)

type strategy interface {
	executeFor(numGR, N int, loopBody func(i, grID int))
}

func defaultStrategy() strategy {
	return new(contiguousBlocksStrategy)
}
