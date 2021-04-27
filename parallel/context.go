package parallel

import "context"

type Context struct {
	context.Context

	loopIndex   int
	goroutineID int
}

func (ctx Context) Index() int {
	return ctx.loopIndex
}

func (ctx Context) GoroutineID() int {
	return ctx.goroutineID
}

func makeParallelContext(ctx context.Context, loopIndex, goroutineID int) *Context {
	return &Context{
		Context:     ctx,
		loopIndex:   loopIndex,
		goroutineID: goroutineID,
	}
}
