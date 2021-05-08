package parallel_test

import (
	"fmt"

	"github.com/dgravesa/go-parallel/parallel"
)

// IncrementNumGRsStrategy implements parallel.Strategy and creates index generators
// corresponding to:
//
//		for i := grID; i < N; i += numGR {
//			// loop body
//		}
//
type IncrementNumGRsStrategy struct{}

func (s *IncrementNumGRsStrategy) IndexGenerator(numGR, grID, _ int) parallel.IndexGenerator {
	return &IncrementNumGRsIndexGenerator{
		increment: numGR,
		nextValue: grID,
	}
}

// IncrementNumGRsIndexGenerator implements parallel.IndexGenerator and returns a sequence of
// indices such that the first index is the goroutine ID and subsequent indices are incremented
// by the total number of goroutines.
type IncrementNumGRsIndexGenerator struct {
	increment int
	nextValue int
}

func (g *IncrementNumGRsIndexGenerator) Next() int {
	thisValue := g.nextValue
	g.nextValue += g.increment
	return thisValue
}

func Example_customStrategy() {
	N := 15
	outputs := make([]int, N)

	strategy := new(IncrementNumGRsStrategy)

	// fill output array with grID used on work index
	parallel.WithCustomStrategy(strategy).WithNumGoroutines(4).For(N, func(i, grID int) {
		outputs[i] = grID
	})

	fmt.Println(outputs)

	// Output: [0 1 2 3 0 1 2 3 0 1 2 3 0 1 2]
}
