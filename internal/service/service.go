package service

import (
	"context"
	"errors"
	"shortener/internal/codec"
	"shortener/internal/errs"
	"shortener/internal/generator"
	"shortener/internal/repository"
	"strings"
)

const maxRetries = 10

type Service struct {
	repo      repository.Repository
	generator generator.Generator
}

func NewService(repo repository.Repository, generator generator.Generator) *Service {
	return &Service{
		repo:      repo,
		generator: generator,
	}
}

func normalizeURL(url string) string {
	url = strings.ToLower(url)
	url = strings.TrimSpace(url)
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimRight(url, "/")
	return url
}

func (s *Service) Shorten(ctx context.Context, longUrl string) (string, error) {
	longUrl = normalizeURL(longUrl)

	existing, err := s.repo.GetShort(ctx, longUrl)
	if err == nil {
		return existing, nil
	}

	for range maxRetries {
		id, err := s.generator.Next(ctx)
		if err != nil {
			return "", err
		}

		shortUrl := codec.Encode(id)

		err = s.repo.Save(ctx, shortUrl, longUrl)
		if err == nil {
			return shortUrl, nil
		}

		if !errors.Is(err, errs.ErrShortLinkAlreadyExists) {
			return "", err
		}
	}

	return "", errs.ErrMaxRetriesExceeded
}

func (s *Service) GetOriginal(ctx context.Context, shortUrl string) (string, error) {

	longUrl, err := s.repo.GetLong(ctx, shortUrl)
	if err != nil {
		return "", err
	}
	longUrl = "https://" + longUrl

	return longUrl, nil
}
