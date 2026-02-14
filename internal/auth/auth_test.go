package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPasswordHashing(t *testing.T) {
	password := "testpassword123"
	
	// Test hashing
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	
	// Test valid password
	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("Failed to check password: %v", err)
	}
	if !match {
		t.Error("Password should match hash")
	}
	
	// Test invalid password
	match, err = CheckPasswordHash("wrongpassword", hash)
	if err != nil {
		t.Fatalf("Failed to check password: %v", err)
	}
	if match {
		t.Error("Wrong password should not match hash")
	}
}

func TestJWTCreationAndValidation(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"
	expiresIn := time.Hour
	
	// Create JWT
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}
	
	// Validate JWT
	parsedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}
	
	if parsedUserID != userID {
		t.Errorf("Expected user ID %v, got %v", userID, parsedUserID)
	}
}

func TestJWTExpiration(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"
	expiresIn := -time.Hour // Already expired
	
	// Create expired JWT
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}
	
	// Try to validate expired JWT
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
}

func TestJWTWrongSecret(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"
	wrongSecret := "wrong-secret-key"
	expiresIn := time.Hour
	
	// Create JWT with one secret
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}
	
	// Try to validate with wrong secret
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("Expected error for wrong secret, got nil")
	}
}

