package parallel

// StrategyType identifies a parallel strategy.
// Different parallel strategies vary on how work items are distributed among goroutines.
// Currently, StrategyContiguousBlocks, StrategyAtomicCounter, and StrategyUseDefaults are the
// accepted values.
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

	// StrategyUseDefaults may be specified on WithStrategy() calls to set the executor to use the
	// default strategies for both For() and ForWithContext().
	StrategyUseDefaults = StrategyType(-1)
)

// Strategy defines a high level interface needed to create a custom index allocation strategy.
// Implementing a custom Strategy is an advanced feature. For most use cases, the strategies
// defines within this package should suffice. Custom Strategy implementations should include an
// IndexGenerator() method, which creates an IndexGenerator instance given the total number of
// goroutines, the ID of the goroutine that the generator pertains to, and the total number of
// work items. An separate IndexGenerator is created for each goroutine. It is the responsibility
// of the custom Strategy implementer to ensure that indices returned by Next() do not overlap
// among the separate IndexGenerator instances.
type Strategy interface {
	IndexGenerator(numGR, grID, N int) IndexGenerator
}

// IndexGenerator defines an interface for individual goroutines to retrieve their work indices.
// IndexGenerator instances should only be created via a corresponding Strategy calling its
// IndexGenerator() method.
// The Next() method of IndexGenerator gives a goroutine its indices to work until all indices
// have been worked, which is specified by Next() returning a value >= total number of loop
// iterations, N.
type IndexGenerator interface {
	Next() int
}
