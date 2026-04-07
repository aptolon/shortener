package generator

import (
	"context"
	"sync/atomic"
)

type Counter struct {
	count atomic.Uint64
}

func (c *Counter) Next(ctx context.Context) (uint64, error) {
	return uint64(c.count.Add(1)), nil
}
