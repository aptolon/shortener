package service

import (
	"context"
	"shortener/internal/codec"
	"shortener/internal/generator"
	"shortener/internal/repository"
)

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

func (s *Service) Shorten(ctx context.Context, longUrl string) (string, error) {
	existing, err := s.repo.GetShort(ctx, longUrl)
	if err == nil {
		return existing, nil
	}

	id, err := s.generator.Next(ctx)
	if err != nil {
		return "", err
	}

	shortUrl := codec.Encode(id)

	if err := s.repo.Save(ctx, shortUrl, longUrl); err != nil {
		return "", err
	}

	return shortUrl, nil

}
