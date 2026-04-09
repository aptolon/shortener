package repository

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestMemoryRepository_SaveAndGet(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	short := "1234567890"
	long := "https://finance.ozon.ru"

	id := uint64(1)
	if err := repo.Save(ctx, id, short, long); err != nil {
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

func TestMemoryRepository_GetNonExistent(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	_, err := repo.GetLong(ctx, "1234567890")
	if err == nil {
		t.Fatal("expected error for non-existent short URL, got nil")
	}

	_, err = repo.GetShort(ctx, "https://finance.ozon.ru")
	if err == nil {
		t.Fatal("expected error for non-existent long URL, got nil")
	}
}

func TestMemoryRepository_SaveDuplicate(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	short := "1234567890"
	long := "https://finance.ozon.ru"

	id := uint64(1)
	if err := repo.Save(ctx, id, short, long); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err := repo.Save(ctx, id+1, short, "https://different.url")
	if err == nil {
		t.Fatal("expected error for duplicate short URL, got nil")
	}

	err = repo.Save(ctx, id+2, "0987654321", long)
	if err == nil {
		t.Fatal("expected error for duplicate long URL, got nil")
	}
}

func TestMemoryRepository_SaveAndGet_Concurrent(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	const numGoroutines = 100

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		short := fmt.Sprintf("1%09d", i)
		long := fmt.Sprintf("https://example.com/%d", i)

		if err := repo.Save(ctx, uint64(i), short, long); err != nil {
			t.Fatalf("pre-save failed: %v", err)
		}
	}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(2)

		go func(i int) {
			defer wg.Done()

			short := fmt.Sprintf("%010d", i)
			long := fmt.Sprintf("https://new.com/%d", i)

			if err := repo.Save(ctx, uint64(1000+i), short, long); err != nil {
				t.Errorf("save failed: %v", err)
			}
		}(i)

		go func(i int) {
			defer wg.Done()

			short := fmt.Sprintf("1%09d", i)
			expectedLong := fmt.Sprintf("https://example.com/%d", i)

			gotLong, err := repo.GetLong(ctx, short)
			if err != nil {
				t.Errorf("get failed: %v", err)
				return
			}

			if gotLong != expectedLong {
				t.Errorf("expected %s, got %s", expectedLong, gotLong)
			}
		}(i)
	}

	wg.Wait()
}
