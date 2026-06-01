package auth

import (
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken("test@example.com", "secret", time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestGenerateToken_EmptySecret(t *testing.T) {
	_, err := GenerateToken("test@example.com", "", time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateToken_Valid(t *testing.T) {
	secret := "my-secret"
	token, err := GenerateToken("user@test.com", secret, time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims.Sub != "user@test.com" {
		t.Errorf("expected sub user@test.com, got %s", claims.Sub)
	}
	if claims.PreferredUsername != "user@test.com" {
		t.Errorf("expected preferred_username user@test.com, got %s", claims.PreferredUsername)
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	token, err := GenerateToken("user@test.com", "correct-secret", time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	_, err = ValidateToken(token, "wrong-secret")
	if err == nil {
		t.Fatal("expected error for wrong secret")
	}
}

func TestValidateToken_Expired(t *testing.T) {
	secret := "secret"
	token, err := GenerateToken("user@test.com", secret, -time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	_, err = ValidateToken(token, secret)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestValidateToken_InvalidString(t *testing.T) {
	_, err := ValidateToken("not-a-token", "secret")
	if err == nil {
		t.Fatal("expected error for invalid token string")
	}
}

func TestValidateToken_EmptyString(t *testing.T) {
	_, err := ValidateToken("", "secret")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestValidateToken_DifferentSigningMethod(t *testing.T) {
	secret := "secret"
	token, err := GenerateToken("user@test.com", secret, time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims.Sub != "user@test.com" {
		t.Errorf("expected sub user@test.com, got %s", claims.Sub)
	}
}
