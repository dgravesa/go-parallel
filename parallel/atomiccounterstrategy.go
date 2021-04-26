package parallel

import (
	"context"
	"sync"
	"sync/atomic"
)

type atomicCounterStrategy struct{}

func (atomicCounterStrategy) executeFor(numGR, N int, loopBody func(i, grID int)) {
	var wg sync.WaitGroup
	wg.Add(numGR)

	// counter and fetcher
	counter := int64(-1)
	fetchIndex := func() int {
		return int(atomic.AddInt64(&counter, 1))
	}

	for grID := 0; grID < numGR; grID++ {
		go func(grID int) {
			defer wg.Done()
			// fetch work indices until work is complete
			for i := fetchIndex(); i < N; i = fetchIndex() {
				loopBody(i, grID)
			}
		}(grID)
	}

	wg.Wait()
}

func (atomicCounterStrategy) executeForWithContext(ctx context.Context, numGR, N int,
	loopBody func(ctx context.Context, i, grID int)) error {

	var wg sync.WaitGroup
	wg.Add(numGR)

	// counter and fetcher
	counter := int64(-1)
	fetchIndex := func() int {
		return int(atomic.AddInt64(&counter, 1))
	}

	for grID := 0; grID < numGR; grID++ {
		go func(grID int) {
			defer wg.Done()
			// fetch work indices until work is complete
			for i := fetchIndex(); i < N; i = fetchIndex() {
				select {
				case <-ctx.Done():
					return
				default:
					loopBody(ctx, i, grID)
				}
			}
		}(grID)
	}

	wg.Wait()

	return ctx.Err()
}
