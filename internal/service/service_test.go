package service

import (
	"context"
	"shortener/internal/generator"
	"shortener/internal/repository"
	"testing"
)

func TestService_Shorten_ReturnsExisting(t *testing.T) {
	repo := repository.NewMemoryRepository()
	gen := generator.NewMemoryGenerator()
	svc := NewService(repo, gen)

	ctx := context.Background()
	url := "https://finance.ozon.ru"

	first, err := svc.Shorten(ctx, url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second, err := svc.Shorten(ctx, url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first != second {
		t.Fatalf("expected same short url, got %s != %s", first, second)
	}
}

func TestService_NormalizeURL(t *testing.T) {

	repo := repository.NewMemoryRepository()
	gen := generator.NewMemoryGenerator()
	svc := NewService(repo, gen)

	ctx := context.Background()

	firstUrl := "finance.ozon.ru"
	secondUrl := "https://fiNanCe.ozon.ru//"

	first, err := svc.Shorten(ctx, firstUrl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second, err := svc.Shorten(ctx, secondUrl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first != second {
		t.Fatalf("expected same short url, got %s != %s", first, second)
	}

}

func TestService_GetOriginal(t *testing.T) {
	repo := repository.NewMemoryRepository()
	gen := generator.NewMemoryGenerator()
	svc := NewService(repo, gen)

	ctx := context.Background()
	url := "https://finance.ozon.ru"

	short, err := svc.Shorten(ctx, url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := svc.GetOriginal(ctx, short)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != url {
		t.Fatalf("expected %s, got %s", url, got)
	}
}
