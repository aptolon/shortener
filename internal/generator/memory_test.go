package generator_test

import (
	"context"
	"shortener/internal/generator"
	"sync"
	"testing"
)

func TestMemoryGenerator_Next_Concurrent(t *testing.T) {
	ctx := context.Background()

	gen := generator.NewMemoryGenerator()

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
