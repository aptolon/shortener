package generator

import "context"

type IDGenerator interface {
	Next(ctx context.Context) (uint64, error)
}
