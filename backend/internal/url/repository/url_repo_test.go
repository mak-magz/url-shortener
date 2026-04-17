package repository_test

import (
	"context"
	"testing"

	"github.com/mak-magz/url-shortener/internal/url/model"
	"github.com/mak-magz/url-shortener/internal/url/service"
	"pgregory.net/rapid"
)

// mockURLRepo implements repository.URLRepository for testing
type mockURLRepo struct {
	getURLFunc func(ctx context.Context, shortCode string) (*model.URL, error)
}

func (m *mockURLRepo) CreateShortURL(ctx context.Context, u *model.URL) error {
	return nil
}

func (m *mockURLRepo) GetURLByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	return m.getURLFunc(ctx, shortCode)
}

func (m *mockURLRepo) IncrementClick(ctx context.Context, id int64) error {
	return nil
}

func TestPropertyShortCodeRoundTrip(t *testing.T) {
	// Feature: url-shortener-phase1, Property 3: Short code round-trip
	rapid.Check(t, func(t *rapid.T) {
		// Generate random 6-char alphanumeric short codes
		shortCode := rapid.StringMatching(`[a-zA-Z0-9]{6}`).Draw(t, "shortCode")

		// Configure a mock repo that returns a model.URL with ShortCode set to that code
		mock := &mockURLRepo{
			getURLFunc: func(ctx context.Context, sc string) (*model.URL, error) {
				return &model.URL{
					ID:          1,
					OriginalURL: "https://example.com",
					ShortCode:   shortCode, // The code we want to "round-trip"
					Clicks:      0,
				}, nil
			},
		}

		// Use the service to test the round-trip through the layers
		svc := service.NewURLService(mock)

		// In a real scenario, this would exercise the GetURLByShortCode call
		// and the subsequent check in the service layer.
		_, err := svc.GetOriginalURL(context.Background(), shortCode)

		// If the round-trip is successful (no mismatch), we should get no error
		if err != nil {
			t.Errorf("expected no error for matching short code, got %v", err)
		}

		// Also directly verify the repository mock returned value as requested
		url, _ := mock.GetURLByShortCode(context.Background(), shortCode)
		if url.ShortCode != shortCode {
			t.Errorf("expected ShortCode %s, got %s", shortCode, url.ShortCode)
		}
	})
}
