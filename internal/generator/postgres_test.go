package generator_test

import (
	"context"
	"os"
	"shortener/internal/generator"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresGenerator_Next_Concurrent(t *testing.T) {
	ctx := context.Background()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Fatal("TEST_DATABASE_URL is not set")
	}

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to create db pool: %v", err)
	}
	defer db.Close()

	gen := generator.NewPostgresGenerator(db)

	const numGoroutines = 10
	const numIDsPerGoroutine = 10
	total := numGoroutines * numIDsPerGoroutine

	results := make(chan uint64, total)
	errors := make(chan error, total)

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := 0; j < numIDsPerGoroutine; j++ {
				id, err := gen.Next(ctx)
				if err != nil {
					errors <- err
					return
				}
				results <- id
			}
		}()
	}

	wg.Wait()
	close(results)
	close(errors)

	for err := range errors {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	seen := make(map[uint64]struct{})

	for id := range results {
		if _, exists := seen[id]; exists {
			t.Fatalf("duplicate id generated: %d", id)
		}
		seen[id] = struct{}{}
	}

	if len(seen) != total {
		t.Fatalf("expected %d ids, got %d", total, len(seen))
	}
}
