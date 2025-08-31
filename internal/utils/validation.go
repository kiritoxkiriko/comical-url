package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// ValidateURL validates if the given string is a valid URL
func ValidateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// Add protocol if missing
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	// Parse URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}

	// Check if host is present
	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a valid host")
	}

	// Check if scheme is http or https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL must use http or https protocol")
	}

	// Additional validation for malformed URLs
	if strings.Contains(parsedURL.Host, " ") {
		return fmt.Errorf("invalid URL: host cannot contain spaces")
	}

	return nil
}

// ValidateCustomKey validates custom short key format
func ValidateCustomKey(key string) error {
	if key == "" {
		return nil // Empty is allowed
	}

	// Check length (3-20 characters)
	if len(key) < 3 || len(key) > 20 {
		return fmt.Errorf("custom key must be between 3 and 20 characters")
	}

	// Check format: alphanumeric, hyphens, underscores only
	matched, err := regexp.MatchString("^[a-zA-Z0-9_-]+$", key)
	if err != nil {
		return fmt.Errorf("error validating custom key: %v", err)
	}

	if !matched {
		return fmt.Errorf("custom key can only contain letters, numbers, hyphens, and underscores")
	}

	// Reserved words check
	reservedWords := []string{"api", "admin", "www", "app", "help", "about", "contact", "terms", "privacy"}
	lowerKey := strings.ToLower(key)
	for _, reserved := range reservedWords {
		if lowerKey == reserved {
			return fmt.Errorf("'%s' is a reserved word and cannot be used as custom key", key)
		}
	}

	return nil
}

// NormalizeURL normalizes URL by adding protocol if missing
func NormalizeURL(rawURL string) string {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return "https://" + rawURL
	}
	return rawURL
}
