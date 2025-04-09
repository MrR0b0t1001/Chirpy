package auth

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := os.Getenv("JWT_SECRET")
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("Expected no error while creating JWT, got %v", err)
	}
	if token == "" {
		t.Errorf("Expected token to be non-empty, got empty string")
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := os.Getenv("JWT_SECRET")
	expiresIn := time.Hour * 1

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("Expected no error while creating JWT, got %v", err)
	}

	validUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Errorf("Expected no error while validating JWT, got %v", err)
	}
	if validUserID != userID {
		t.Errorf("Expected userID %v, got %v", userID, validUserID)
	}
}

func TestValidateExpiredJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := os.Getenv("JWT_SECRET")
	expiresIn := -time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("Expected no error while creating JWT, got %v", err)
	}

	_, err = ValidateJWT(token, tokenSecret)

	if err == nil {
		t.Errorf("Expected error while validating expired JWT, got nil")
	}
}

func TestValidateJWTWithWrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secretKey"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("Expected no error while creating JWT, got %v", err)
	}

	incorrectSecret := "wrongSecret"
	_, err = ValidateJWT(token, incorrectSecret)

	if err == nil {
		t.Errorf("Expected error while validating with incorrect secret, got nil")
	}
}

func TestValidateMalformedJWT(t *testing.T) {
	invalidToken := "invalidTokenString"
	tokenSecret := os.Getenv("JWT_SECRET")

	_, err := ValidateJWT(invalidToken, tokenSecret)

	if err == nil {
		t.Errorf("Expected error while validating malformed JWT, got nil")
	}
}
