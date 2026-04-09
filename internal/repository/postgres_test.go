package repository

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresRepository_SaveAndGet(t *testing.T) {
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Fatal("DATABASE_URL is required")
	}

	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}
	defer db.Close()
	repo := NewPostgresRepository(db)

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
