package utils

import "testing"

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid https URL",
			url:     "https://example.com",
			wantErr: false,
		},
		{
			name:    "valid http URL",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "URL without protocol",
			url:     "example.com",
			wantErr: false,
		},
		{
			name:    "URL with path",
			url:     "https://example.com/path/to/page",
			wantErr: false,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "invalid URL format",
			url:     "ht tp://invalid url.com",
			wantErr: true,
		},
		{
			name:    "URL without host",
			url:     "https://",
			wantErr: true,
		},
		{
			name:    "invalid protocol scheme",
			url:     "javascript:alert('xss')",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCustomKey(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name:    "valid key",
			key:     "mykey",
			wantErr: false,
		},
		{
			name:    "valid key with numbers",
			key:     "key123",
			wantErr: false,
		},
		{
			name:    "valid key with hyphen",
			key:     "my-key",
			wantErr: false,
		},
		{
			name:    "valid key with underscore",
			key:     "my_key",
			wantErr: false,
		},
		{
			name:    "empty key (allowed)",
			key:     "",
			wantErr: false,
		},
		{
			name:    "too short key",
			key:     "ab",
			wantErr: true,
		},
		{
			name:    "too long key",
			key:     "this-is-a-very-long-key-name",
			wantErr: true,
		},
		{
			name:    "key with special characters",
			key:     "key@123",
			wantErr: true,
		},
		{
			name:    "key with spaces",
			key:     "my key",
			wantErr: true,
		},
		{
			name:    "reserved word",
			key:     "api",
			wantErr: true,
		},
		{
			name:    "reserved word case insensitive",
			key:     "API",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCustomKey(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCustomKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "URL with https",
			input:    "https://example.com",
			expected: "https://example.com",
		},
		{
			name:     "URL with http",
			input:    "http://example.com",
			expected: "http://example.com",
		},
		{
			name:     "URL without protocol",
			input:    "example.com",
			expected: "https://example.com",
		},
		{
			name:     "URL with path",
			input:    "example.com/path",
			expected: "https://example.com/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeURL(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeURL() = %v, want %v", result, tt.expected)
			}
		})
	}
}