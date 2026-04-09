package repository_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	"shortener/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newTestRepo(t *testing.T) (*repository.PostgresRepository, context.Context, *pgxpool.Pool) {
	t.Helper()

	ctx := context.Background()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Fatal("TEST_DATABASE_URL is not set")
	}

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	_, err = db.Exec(ctx, `TRUNCATE TABLE urls RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("failed to truncate table: %v", err)
	}

	repo := repository.NewPostgresRepository(db)
	return repo, ctx, db
}
func TestPostgresRepository_SaveAndGet(t *testing.T) {
	repo, ctx, db := newTestRepo(t)
	defer db.Close()

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

func TestPostgresRepository_GetNonExistent(t *testing.T) {
	repo, ctx, db := newTestRepo(t)
	defer db.Close()

	_, err := repo.GetLong(ctx, "1234567890")
	if err == nil {
		t.Fatal("expected error for non-existent short URL, got nil")
	}

	_, err = repo.GetShort(ctx, "https://finance.ozon.ru")
	if err == nil {
		t.Fatal("expected error for non-existent long URL, got nil")
	}

}
func TestPostgresRepository_SaveDuplicate(t *testing.T) {
	repo, ctx, db := newTestRepo(t)
	defer db.Close()

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

func TestPostgresRepository_SaveAndGet_Concurrent(t *testing.T) {
	repo, ctx, db := newTestRepo(t)
	defer db.Close()

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
