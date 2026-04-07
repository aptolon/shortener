package generator

import (
	"context"
	"sync/atomic"
)

type MemoryGenerator struct {
	counter atomic.Uint64
}

func NewMemoryGenerator() *MemoryGenerator {
	return &MemoryGenerator{}
}

func (g *MemoryGenerator) Next(ctx context.Context) (uint64, error) {
	return g.counter.Add(1), nil
}
