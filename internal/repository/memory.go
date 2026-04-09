package repository

import (
	"context"
	"sync"

	"shortener/internal/errs"
)

type MemoryRepository struct {
	mu          sync.RWMutex
	shortToLong map[string]string
	longToShort map[string]string
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		shortToLong: make(map[string]string),
		longToShort: make(map[string]string),
	}
}

func (r *MemoryRepository) Save(ctx context.Context, id uint64, shortUrl string, longUrl string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.shortToLong[shortUrl]; ok {
		return errs.ErrShortLinkAlreadyExists
	}
	if _, ok := r.longToShort[longUrl]; ok {
		return errs.ErrLongLinkAlreadyExists
	}
	r.shortToLong[shortUrl] = longUrl
	r.longToShort[longUrl] = shortUrl

	return nil
}

func (r *MemoryRepository) GetLong(ctx context.Context, shortUrl string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if longUrl, ok := r.shortToLong[shortUrl]; ok {
		return longUrl, nil
	}
	return "", errs.ErrNotFound
}

func (r *MemoryRepository) GetShort(ctx context.Context, longUrl string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if shortUrl, ok := r.longToShort[longUrl]; ok {
		return shortUrl, nil
	}
	return "", errs.ErrNotFound
}
