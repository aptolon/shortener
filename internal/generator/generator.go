package generator

import "context"

type Generator interface {
	Next(ctx context.Context) (uint64, error)
}
