package service

import "shortener/internal/repository"

type Service struct {
	repo repository.Repository
}

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func NewService(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}
