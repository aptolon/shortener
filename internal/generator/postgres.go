package generator

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresGenerator struct {
	db *pgxpool.Pool
}

func NewPostgresGenerator(db *pgxpool.Pool) *PostgresGenerator {
	return &PostgresGenerator{
		db: db,
	}
}

func (g *PostgresGenerator) Next(ctx context.Context) (uint64, error) {
	var id uint64

	err := g.db.QueryRow(ctx, `SELECT nextval('links_id_seq')`).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
