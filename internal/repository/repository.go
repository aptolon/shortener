package repository

import "context"

type Repository interface {
	Save(ctx context.Context, id uint64, shortUrl string, longUrl string) error
	GetLong(ctx context.Context, shortUrl string) (string, error)
	GetShort(ctx context.Context, longUrl string) (string, error)
}
