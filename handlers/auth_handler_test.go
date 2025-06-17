package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestCreateTokenRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request CreateTokenRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: CreateTokenRequest{
				Name: "Test Token",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			request: CreateTokenRequest{
				Name: "",
			},
			wantErr: true,
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
			req := httptest.NewRequest("POST", "/api/auth/tokens", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			// Test JSON binding
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			var testReq CreateTokenRequest
			err = c.ShouldBindJSON(&testReq)

			if tt.wantErr {
				// Should fail validation for missing required field
				if tt.request.Name == "" && err == nil {
					t.Error("Expected validation error for missing name")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected validation error: %v", err)
				}
				if testReq.Name != tt.request.Name {
					t.Errorf("Name mismatch: got %v, want %v", testReq.Name, tt.request.Name)
				}
			}
		})
	}
}

func TestCreateTokenResponse_JSON(t *testing.T) {
	response := CreateTokenResponse{
		Token: "test-token-123",
		Name:  "Test Token",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled CreateTokenResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify fields
	if unmarshaled.Token != response.Token {
		t.Errorf("Token mismatch: got %v, want %v", unmarshaled.Token, response.Token)
	}
	if unmarshaled.Name != response.Name {
		t.Errorf("Name mismatch: got %v, want %v", unmarshaled.Name, response.Name)
	}
}

func TestAuthHandler_Creation(t *testing.T) {
	handler := NewAuthHandler()
	if handler == nil {
		t.Error("NewAuthHandler() returned nil")
	}
}

func TestUUID_Generation(t *testing.T) {
	// Test UUID generation
	token1 := uuid.New().String()
	token2 := uuid.New().String()

	// Tokens should be different
	if token1 == token2 {
		t.Errorf("Generated UUIDs should be different: %s == %s", token1, token2)
	}

	// Test UUID format (36 characters including hyphens)
	if len(token1) != 36 {
		t.Errorf("Expected UUID length 36, got %d", len(token1))
	}

	// Test UUID parsing
	_, err := uuid.Parse(token1)
	if err != nil {
		t.Errorf("Generated token is not a valid UUID: %v", err)
	}
}