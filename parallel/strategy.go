package parallel

// StrategyType identifies a parallel strategy.
// Different parallel strategies vary on how work items are executed across goroutines.
// The available strategy types are defined as constants and follow the naming convention
// Strategy*.
type StrategyType int

const (
	// StrategyContiguousBlocks refers to a strategy that preassigns an equal number of iterations
	// per goroutine as contiguous blocks.
	// The contiguous blocks strategy works well for most parallelizable loops.
	// The atomic counter strategy may be faster when loop iterations vary in execution time or
	// when the execution time per iteration is greater than 1 microsecond.
	StrategyContiguousBlocks = StrategyType(iota)

	// StrategyAtomicCounter refers to a strategy that uses an atomic counter so goroutines can
	// fetch additional work items when they are ready as opposed to preassigning iterations
	// to each goroutine.
	// This strategy may be faster than the contiguous blocks strategy for parallelizable loops
	// with greater time variation across iterations.
	StrategyAtomicCounter = StrategyType(iota)
)

type strategy interface {
	executeFor(numGR, N int, loopBody func(i, grID int))
}

func defaultStrategy() strategy {
	return new(contiguousBlocksStrategy)
}
