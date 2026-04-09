package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"shortener/internal/generator"
	"shortener/internal/repository"
)

func BuildStorage(
	ctx context.Context,
	storage string,
	databaseURL string,
) (repository.Repository, generator.Generator, func(), error) {
	switch storage {
	case "memory":
		return repository.NewMemoryRepository(),
			generator.NewMemoryGenerator(),
			func() {},
			nil

	case "postgres":
		if databaseURL == "" {
			return nil, nil, nil, fmt.Errorf("DATABASE_URL is required")
		}

		db, err := pgxpool.New(ctx, databaseURL)
		if err != nil {
			return nil, nil, nil, err
		}

		if err := db.Ping(ctx); err != nil {
			db.Close()
			return nil, nil, nil, err
		}

		return repository.NewPostgresRepository(db),
			generator.NewPostgresGenerator(db),
			func() { db.Close() },
			nil

	default:
		return nil, nil, nil, fmt.Errorf("unknown storage: %s", storage)
	}
}
