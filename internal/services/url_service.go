package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"

	"shorturl/internal/config"
	"shorturl/internal/models"
	"shorturl/internal/utils"
)

type URLService struct{}

func NewURLService() *URLService {
	return &URLService{}
}

func (s *URLService) GenerateShortKey() string {
	key, err := gonanoid.New(6)
	if err != nil {
		// Fallback to a simple implementation if nanoid fails
		return "fallbk"
	}
	return key
}

func (s *URLService) CreateShortURL(longURL, customKey, passkey, expiresIn string) (*models.URL, error) {
	// Validate URL
	if err := utils.ValidateURL(longURL); err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	// Normalize URL
	longURL = utils.NormalizeURL(longURL)

	// Validate custom key if provided
	if err := utils.ValidateCustomKey(customKey); err != nil {
		return nil, fmt.Errorf("invalid custom key: %v", err)
	}

	var shortKey string
	if customKey != "" {
		// Check if custom key already exists
		var existingURL models.URL
		if err := config.DB.Where("short_key = ?", customKey).First(&existingURL).Error; err == nil {
			return nil, errors.New("custom key already exists")
		}
		shortKey = customKey
	} else {
		// Generate unique short key using nanoid
		for {
			shortKey = s.GenerateShortKey()
			var existingURL models.URL
			if err := config.DB.Where("short_key = ?", shortKey).First(&existingURL).Error; err != nil {
				break // Key doesn't exist, we can use it
			}
		}
	}

	var passkeyHash string
	if passkey != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(passkey), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash passkey: %v", err)
		}
		passkeyHash = string(hash)
	}

	var expiresAt *time.Time
	if expiresIn != "" {
		duration, err := time.ParseDuration(expiresIn)
		if err != nil {
			return nil, fmt.Errorf("invalid expires_in format: %v", err)
		}
		expiry := time.Now().Add(duration)
		expiresAt = &expiry
	} else {
		// Use default expiration from config or fallback
		duration, err := time.ParseDuration("720h") // 30 days default
		if err == nil {
			expiry := time.Now().Add(duration)
			expiresAt = &expiry
		}
	}

	url := &models.URL{
		ShortKey:    shortKey,
		LongURL:     longURL,
		ExpiresAt:   expiresAt,
		PasskeyHash: passkeyHash,
		IsActive:    true,
	}

	if err := config.DB.Create(url).Error; err != nil {
		return nil, fmt.Errorf("failed to create short URL: %v", err)
	}

	// Cache in Redis for faster access
	ctx := context.Background()
	config.Redis.Set(ctx, "url:"+shortKey, longURL, time.Hour*24*7) // Cache for 7 days

	return url, nil
}

func (s *URLService) GetLongURL(shortKey, passkey string) (string, error) {
	if shortKey == "" {
		return "", errors.New("short key is required")
	}

	// Try Redis cache first
	ctx := context.Background()
	cachedURL, err := config.Redis.Get(ctx, "url:"+shortKey).Result()
	if err == nil && cachedURL != "" {
		// Still need to check passkey and update clicks in DB
		var url models.URL
		if err := config.DB.Where("short_key = ? AND is_active = ?", shortKey, true).First(&url).Error; err != nil {
			return "", errors.New("URL not found")
		}

		if err := s.validateURL(&url, passkey); err != nil {
			return "", err
		}

		// Update clicks
		config.DB.Model(&url).Update("clicks", url.Clicks+1)

		return cachedURL, nil
	}

	// Fallback to database
	var url models.URL
	if err := config.DB.Where("short_key = ? AND is_active = ?", shortKey, true).First(&url).Error; err != nil {
		return "", errors.New("URL not found")
	}

	if err := s.validateURL(&url, passkey); err != nil {
		return "", err
	}

	// Update clicks
	config.DB.Model(&url).Update("clicks", url.Clicks+1)

	// Cache the result
	config.Redis.Set(ctx, "url:"+shortKey, url.LongURL, time.Hour*24*7)

	return url.LongURL, nil
}

func (s *URLService) validateURL(url *models.URL, passkey string) error {
	// Check if URL is expired
	if url.ExpiresAt != nil && time.Now().After(*url.ExpiresAt) {
		return errors.New("URL has expired")
	}

	// Check passkey if required
	if url.PasskeyHash != "" {
		if passkey == "" {
			return errors.New("passkey required")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(url.PasskeyHash), []byte(passkey)); err != nil {
			return errors.New("invalid passkey")
		}
	}

	return nil
}

func (s *URLService) RevokeURL(shortKey string) error {
	result := config.DB.Model(&models.URL{}).Where("short_key = ?", shortKey).Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("URL not found")
	}

	// Remove from cache
	ctx := context.Background()
	config.Redis.Del(ctx, "url:"+shortKey)

	return nil
}

func (s *URLService) AutoRevokeExpiredURLs() error {
	now := time.Now()
	result := config.DB.Model(&models.URL{}).
		Where("expires_at IS NOT NULL AND expires_at < ? AND is_active = ?", now, true).
		Update("is_active", false)

	return result.Error
}
