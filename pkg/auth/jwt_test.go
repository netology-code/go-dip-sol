package auth

import (
	"testing"
	"time"
)

func TestJWTManager_GenerateToken(t *testing.T) {
	manager := NewJWTManager("test-secret", 24)

	userID := 1
	email := "test@example.com"
	username := "testuser"

	token, expiresAt, err := manager.GenerateToken(userID, email, username)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Error("expected non-empty token")
	}

	if expiresAt.Before(time.Now()) {
		t.Error("expected expiresAt in the future")
	}
}

func TestJWTManager_ValidateToken_Success(t *testing.T) {
	manager := NewJWTManager("test-secret", 24)

	userID := 1
	email := "test@example.com"
	username := "testuser"

	token, _, err := manager.GenerateToken(userID, email, username)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected UserID %d, got %d", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("expected Email %s, got %s", email, claims.Email)
	}

	if claims.Username != username {
		t.Errorf("expected Username %s, got %s", username, claims.Username)
	}
}

func TestJWTManager_ValidateToken_InvalidToken(t *testing.T) {
	manager := NewJWTManager("test-secret", 24)

	invalidToken := "invalid.token.here"

	_, err := manager.ValidateToken(invalidToken)
	if err == nil {
		t.Error("expected error for invalid token")
	}
}
