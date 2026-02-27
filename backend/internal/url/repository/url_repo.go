package url

import (
	"context"

	"github.com/mak-magz/url-shortener/internal/url/model"
)

type URLRepository interface {
	CreateShortURL(ctx context.Context, req *model.URL) error
	GetURLByShortCode(ctx context.Context, shortCode string) (*model.URL, error)
	IncrementClick(ctx context.Context, id int64) error
}
