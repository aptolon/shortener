package repository

import (
	"context"
	"testing"
)

func TestMemoryRepository_SaveAndGet(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	short := "1234567890"
	long := "https://finance.ozon.ru"

	if err := repo.Save(ctx, short, long); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotLong, err := repo.GetLong(ctx, short)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotLong != long {
		t.Fatalf("expected %s, got %s", long, gotLong)
	}

	gotShort, err := repo.GetShort(ctx, long)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotShort != short {
		t.Fatalf("expected %s, got %s", short, gotShort)
	}

}
