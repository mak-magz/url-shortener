package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/mak-magz/url-shortener/internal/url/model"
	"github.com/mak-magz/url-shortener/internal/url/service"
)

type MockRepo struct {
}

func (m *MockRepo) CreateShortURL(ctx context.Context, u *model.URL) error {

	u.ID = 1
	u.Clicks = 0
	u.CreatedAt = time.Now()
	return nil
}

func (m *MockRepo) GetURLByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	url := &model.URL{
		ID:          1,
		OriginalURL: "https://www.google.com",
		ShortCode:   "abc123",
		Clicks:      0,
		CreatedAt:   time.Now(),
	}
	return url, nil
}

func (m *MockRepo) IncrementClick(ctx context.Context, id int64) error {
	return nil
}

func TestCreateShortURL(t *testing.T) {
	mockRepo := &MockRepo{}
	service := service.NewURLService(mockRepo)

	req := &model.CreateURLRequest{
		OriginalURL: "https://www.google.com",
	}

	resp, err := service.CreateShortURL(context.Background(), req)
	if err != nil {
		t.Errorf("Error creating short URL: %v", err)
	}

	if resp.OriginalURL != req.OriginalURL {
		t.Errorf("Original URL does not match")
	}

	if resp.ShortCode == "" {
		t.Errorf("Short code is empty")
	}

	if len(resp.ShortCode) != 6 {
		t.Errorf("Short code length is not 6")
	}

	if resp.Clicks != 0 {
		t.Errorf("Clicks should be 0")
	}

	if resp.CreatedAt.IsZero() {
		t.Errorf("Created at is zero")
	}
}

func TestGetOriginalURL(t *testing.T) {
	mockRepo := &MockRepo{}
	service := service.NewURLService(mockRepo)

	resp, err := service.GetOriginalURL(context.Background(), "abc123")
	if err != nil {
		t.Errorf("Error getting original URL: %v", err)
	}

	if resp != "https://www.google.com" {
		t.Errorf("Original URL does not match")
	}
}
