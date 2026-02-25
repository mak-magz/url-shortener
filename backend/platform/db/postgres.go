package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect creates a connection pool to PostgreSQL
func Connect(databaseURL string) *pgxpool.Pool {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse database URL
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v\n", err)
	}

	// Connection pool settings
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Ping the database to verify the connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	log.Println("✅ Database connection established")

	return pool
}

// Migrate runs database migrations
func Migrate(pool *pgxpool.Pool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Run database migrations
	query := `
		CREATE TABLE IF NOT EXISTS urls (
			id 							SERIAL PRIMARY KEY,
			original_url 		TEXT NOT NULL,
			short_code 			VARCHAR(10) NOT NULL,
			clicks 					INT DEFAULT 0,
			created_at 			TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);

		CREATE UNIQUE INDEX IF NOT EXISTS idx_short_code ON urls(short_code);
	`

	_, err := pool.Exec(ctx, query)
	if err != nil {
		log.Fatalf("Unable to run database migrations: %v\n", err)
	}

	log.Println("✅ Database migrations completed")
}
