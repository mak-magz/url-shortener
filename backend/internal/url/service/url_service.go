package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/mak-magz/url-shortener/internal/url/model"
	"github.com/mak-magz/url-shortener/internal/url/repository"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const codeLength = 6

type Service interface {
	CreateShortURL(ctx context.Context, req *model.CreateURLRequest) (*model.CreateURLResponse, error)
	GetOriginalURL(ctx context.Context, shortCode string) (string, error)
}

type URLService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) Service {
	return &URLService{repo: repo}
}

// CreateShortURL implements [Service].
func (u *URLService) CreateShortURL(ctx context.Context, req *model.CreateURLRequest) (*model.CreateURLResponse, error) {
	shortCode := generateShortCode(codeLength)

	newURL := &model.URL{
		OriginalURL: req.OriginalURL,
		ShortCode:   shortCode,
	}

	if err := u.repo.CreateShortURL(ctx, newURL); err != nil {
		return nil, err
	}

	return &model.CreateURLResponse{
		ID:          newURL.ID,
		OriginalURL: newURL.OriginalURL,
		ShortCode:   newURL.ShortCode,
		Clicks:      newURL.Clicks,
		CreatedAt:   newURL.CreatedAt,
	}, nil
}

// GetOriginalURL implements [Service].
func (u *URLService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	url, err := u.repo.GetURLByShortCode(ctx, shortCode)

	if err != nil {
		return "", err
	}

	err = u.repo.IncrementClick(ctx, url.ID)

	if err != nil {
		return "", err
	}

	return url.OriginalURL, nil
}

func generateShortCode(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[r.Intn(len(charset))]
	}
	return string(code)
}
