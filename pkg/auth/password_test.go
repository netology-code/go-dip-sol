package auth

import (
	"testing"
)

func TestHashPassword_Success(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hash == "" {
		t.Error("expected non-empty hash")
	}

	// Hash should start with bcrypt identifier
	if hash[:4] != "$2a$" && hash[:4] != "$2b$" && hash[:4] != "$2y$" {
		t.Error("expected bcrypt hash format")
	}
}

func TestCheckPassword_Success(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("hashing failed: %v", err)
	}

	if !CheckPassword(password, hash) {
		t.Error("expected password to match hash")
	}
}

func TestCheckPassword_Failure_WrongPassword(t *testing.T) {
	password := "testpassword123"
	wrongPassword := "wrongpassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("hashing failed: %v", err)
	}

	if CheckPassword(wrongPassword, hash) {
		t.Error("expected password not to match hash")
	}
}

func TestCheckPassword_Failure_InvalidHash(t *testing.T) {
	password := "testpassword123"
	invalidHash := "invalidhash"

	if CheckPassword(password, invalidHash) {
		t.Error("expected invalid hash to fail")
	}
}
