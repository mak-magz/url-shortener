package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mak-magz/url-shortener/internal/url/model"
)

type URLRepository interface {
	CreateShortURL(ctx context.Context, req *model.URL) error
	GetURLByShortCode(ctx context.Context, shortCode string) (*model.URL, error)
	IncrementClick(ctx context.Context, id int64) error
}

var ErrURLNotFound = errors.New("url not found")

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) URLRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) CreateShortURL(ctx context.Context, u *model.URL) error {
	query := `
			INSERT INTO urls (original_url, short_code) 
			VALUES ($1, $2)
			RETURNING id, clicks, created_at	
		`

	err := r.pool.QueryRow(ctx, query, u.OriginalURL, u.ShortCode).Scan(&u.ID, &u.Clicks, &u.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) GetURLByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	query := `
		SELECT id, original_url, clicks, created_at, updated_at
		FROM urls
		WHERE short_code = $1
	`

	u := &model.URL{}

	err := r.pool.QueryRow(ctx, query, shortCode).Scan(&u.ID, &u.OriginalURL, &u.Clicks, &u.CreatedAt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *PostgresRepository) IncrementClick(ctx context.Context, id int64) error {
	query := `
		UPDATE urls
		SET clicks = clicks + 1
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
