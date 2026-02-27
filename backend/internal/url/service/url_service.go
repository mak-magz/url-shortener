package service

import (
	"context"

	"github.com/mak-magz/url-shortener/internal/url/model"
)

type URLService interface {
	CreateShortURL(ctx context.Context, req *model.CreateURLRequest) (*model.CreateURLResponse, error)
	GetOriginalURL(ctx context.Context, shortCode string) (string, error)
}
