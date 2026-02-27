package model

import "time"

type URL struct {
	ID          int64     `db:"id" json:"id"`
	OriginalURL string    `db:"original_url" json:"originalUrl"`
	ShortCode   string    `db:"short_code" json:"shortCode"`
	Clicks      int64     `db:"clicks" json:"clicks"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}

type CreateURLRequest struct {
	OriginalURL string `json:"originalUrl" binding:"required"`
}

type CreateURLResponse struct {
	ID          int64     `json:"id"`
	OriginalURL string    `json:"originalUrl"`
	ShortCode   string    `json:"shortCode"`
	Clicks      int64     `json:"clicks"`
	CreatedAt   time.Time `json:"createdAt"`
}
