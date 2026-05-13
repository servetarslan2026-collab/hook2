package auth

import (
	"testing"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("testpassword123")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword returned empty hash")
	}
	if hash == "testpassword123" {
		t.Fatal("HashPassword returned plaintext")
	}
}

func TestCheckPassword(t *testing.T) {
	hash, err := HashPassword("testpassword123")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if !CheckPassword("testpassword123", hash) {
		t.Fatal("CheckPassword failed for correct password")
	}

	if CheckPassword("wrongpassword", hash) {
		t.Fatal("CheckPassword succeeded for wrong password")
	}
}

func TestGenerateToken(t *testing.T) {
	svc := NewAuthService("test-secret-key-min-32-chars!!", "24h")

	userID := uuid.New()
	token, err := svc.GenerateToken(userID, "test@example.com")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("GenerateToken returned empty token")
	}

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if claims.UserID != userID {
		t.Fatalf("Expected UserID %s, got %s", userID, claims.UserID)
	}
	if claims.Email != "test@example.com" {
		t.Fatalf("Expected Email 'test@example.com', got '%s'", claims.Email)
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	svc := NewAuthService("test-secret-key-min-32-chars!!", "24h")

	_, err := svc.ValidateToken("invalid.token.here")
	if err == nil {
		t.Fatal("ValidateToken should fail for invalid token")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	svc := NewAuthService("test-secret-key-min-32-chars!!", "24h")

	userID := uuid.New()
	token, err := svc.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("GenerateRefreshToken returned empty token")
	}

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed for refresh token: %v", err)
	}
	if claims.UserID != userID {
		t.Fatalf("Expected UserID %s, got %s", userID, claims.UserID)
	}
}

func TestGenerateAPIKey(t *testing.T) {
	key := GenerateAPIKey()
	if key == "" {
		t.Fatal("GenerateAPIKey returned empty")
	}
	if len(key) < 10 {
		t.Fatalf("API key too short: %s", key)
	}

	// Generate two keys, they should be different
	key2 := GenerateAPIKey()
	if key == key2 {
		t.Fatal("GenerateAPIKey returned duplicate keys")
	}
}
