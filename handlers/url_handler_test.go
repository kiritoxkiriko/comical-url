package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateURLRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request CreateURLRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: CreateURLRequest{
				LongURL:   "https://example.com",
				CustomKey: "mykey",
				Passkey:   "secret",
				ExpiresIn: "24h",
			},
			wantErr: false,
		},
		{
			name: "missing long URL",
			request: CreateURLRequest{
				CustomKey: "mykey",
				Passkey:   "secret",
				ExpiresIn: "24h",
			},
			wantErr: true,
		},
		{
			name: "valid minimal request",
			request: CreateURLRequest{
				LongURL: "https://example.com",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert request to JSON
			jsonData, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}

			// Create HTTP request
			req := httptest.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			// Test JSON binding
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			var testReq CreateURLRequest
			err = c.ShouldBindJSON(&testReq)

			if tt.wantErr {
				// Should fail validation for missing required field
				if tt.request.LongURL == "" && err == nil {
					t.Error("Expected validation error for missing long_url")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected validation error: %v", err)
				}
				if testReq.LongURL != tt.request.LongURL {
					t.Errorf("LongURL mismatch: got %v, want %v", testReq.LongURL, tt.request.LongURL)
				}
			}
		})
	}
}

func TestCreateURLResponse_JSON(t *testing.T) {
	response := CreateURLResponse{
		ShortKey: "abc123",
		ShortURL: "https://short.ly/abc123",
		LongURL:  "https://example.com/very/long/url",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled CreateURLResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify fields
	if unmarshaled.ShortKey != response.ShortKey {
		t.Errorf("ShortKey mismatch: got %v, want %v", unmarshaled.ShortKey, response.ShortKey)
	}
	if unmarshaled.ShortURL != response.ShortURL {
		t.Errorf("ShortURL mismatch: got %v, want %v", unmarshaled.ShortURL, response.ShortURL)
	}
	if unmarshaled.LongURL != response.LongURL {
		t.Errorf("LongURL mismatch: got %v, want %v", unmarshaled.LongURL, response.LongURL)
	}
}

func TestURLHandler_Creation(t *testing.T) {
	handler := NewURLHandler()
	if handler == nil {
		t.Error("NewURLHandler() returned nil")
	}
	if handler.urlService == nil {
		t.Error("URLHandler.urlService is nil")
	}
}