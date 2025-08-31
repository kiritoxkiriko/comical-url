package services

import (
	"testing"
	"time"

	"github.com/matoous/go-nanoid/v2"
	"shorturl/internal/models"
)

func TestGenerateShortKey(t *testing.T) {
	service := NewURLService()

	// Test key generation
	key1 := service.GenerateShortKey()
	key2 := service.GenerateShortKey()

	// Keys should be 6 characters long
	if len(key1) != 6 {
		t.Errorf("Expected key length 6, got %d", len(key1))
	}

	// Keys should be different
	if key1 == key2 {
		t.Errorf("Generated keys should be different: %s == %s", key1, key2)
	}

	// Test with specific alphabet for nanoid
	customKey, err := gonanoid.New(6)
	if err != nil {
		t.Errorf("Failed to generate nanoid: %v", err)
	}

	if len(customKey) != 6 {
		t.Errorf("Expected nanoid length 6, got %d", len(customKey))
	}
}

func TestValidateURL(t *testing.T) {
	service := NewURLService()

	tests := []struct {
		name      string
		url       string
		passkey   string
		expiresAt *time.Time
		wantErr   bool
	}{
		{
			name:      "valid URL not expired",
			url:       "https://example.com",
			passkey:   "",
			expiresAt: nil,
			wantErr:   false,
		},
		{
			name:      "expired URL",
			url:       "https://example.com",
			passkey:   "",
			expiresAt: timePtr(time.Now().Add(-1 * time.Hour)),
			wantErr:   true,
		},
		{
			name:      "future expiry",
			url:       "https://example.com",
			passkey:   "",
			expiresAt: timePtr(time.Now().Add(1 * time.Hour)),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := &models.URL{
				LongURL:   tt.url,
				ExpiresAt: tt.expiresAt,
			}

			err := service.validateURL(url, tt.passkey)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function to create time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration string
		wantErr  bool
	}{
		{
			name:     "seconds",
			duration: "30s",
			wantErr:  false,
		},
		{
			name:     "minutes",
			duration: "5m",
			wantErr:  false,
		},
		{
			name:     "hours",
			duration: "2h",
			wantErr:  false,
		},
		{
			name:     "days (as hours)",
			duration: "168h", // 7 days
			wantErr:  false,
		},
		{
			name:     "invalid format",
			duration: "7d", // days not supported by time.ParseDuration
			wantErr:  true,
		},
		{
			name:     "invalid string",
			duration: "invalid",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := time.ParseDuration(tt.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDuration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
