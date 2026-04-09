package repository

import (
	"context"
	"shortener/internal/errs"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Save(ctx context.Context, id uint64, shortUrl string, longUrl string) error {
	query := `
		INSERT INTO links (id, short_code, long_url)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(ctx, query, id, shortUrl, longUrl)
	if err != nil {
		return errs.ErrShortLinkAlreadyExists
	}
	return err
}

func (r *PostgresRepository) GetLong(ctx context.Context, shortUrl string) (string, error) {
	query := `
		SELECT long_url FROM links WHERE short_code = $1
	`

	var longUrl string
	err := r.db.QueryRow(ctx, query, shortUrl).Scan(&longUrl)
	if err != nil {
		return "", errs.ErrNotFound
	}
	return longUrl, nil
}

func (r *PostgresRepository) GetShort(ctx context.Context, longUrl string) (string, error) {
	query := `
		SELECT short_code FROM links WHERE long_url = $1
	`

	var shortUrl string
	err := r.db.QueryRow(ctx, query, longUrl).Scan(&shortUrl)
	if err != nil {
		return "", errs.ErrNotFound
	}
	return shortUrl, nil
}
