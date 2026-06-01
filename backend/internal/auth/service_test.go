package auth

import (
	"testing"
	"time"
)

func TestAuthenticate_ValidCredentials(t *testing.T) {
	service := NewAuthService("secret", time.Hour)
	token, err := service.Authenticate("tecnico@prefeitura.rio", "painel@2024")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := ValidateToken(token, "secret")
	if err != nil {
		t.Fatalf("failed to validate generated token: %v", err)
	}
	if claims.Sub != "tecnico@prefeitura.rio" {
		t.Errorf("expected sub tecnico@prefeitura.rio, got %s", claims.Sub)
	}
}

func TestAuthenticate_InvalidEmail(t *testing.T) {
	service := NewAuthService("secret", time.Hour)
	_, err := service.Authenticate("wrong@email.com", "painel@2024")
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthenticate_InvalidPassword(t *testing.T) {
	service := NewAuthService("secret", time.Hour)
	_, err := service.Authenticate("tecnico@prefeitura.rio", "wrong-password")
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthenticate_EmptyCredentials(t *testing.T) {
	service := NewAuthService("secret", time.Hour)
	_, err := service.Authenticate("", "")
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthenticate_EmptyEmail(t *testing.T) {
	service := NewAuthService("secret", time.Hour)
	_, err := service.Authenticate("", "painel@2024")
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthenticate_EmptyPassword(t *testing.T) {
	service := NewAuthService("secret", time.Hour)
	_, err := service.Authenticate("tecnico@prefeitura.rio", "")
	if err != ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestServiceValidateToken_Valid(t *testing.T) {
	service := NewAuthService("my-secret", time.Hour)
	token, err := service.Authenticate("tecnico@prefeitura.rio", "painel@2024")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}

	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims.Sub != "tecnico@prefeitura.rio" {
		t.Errorf("expected sub tecnico@prefeitura.rio, got %s", claims.Sub)
	}
}

func TestServiceValidateToken_Invalid(t *testing.T) {
	service := NewAuthService("my-secret", time.Hour)
	_, err := service.ValidateToken("invalid-token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestServiceValidateToken_WrongSecret(t *testing.T) {
	serviceA := NewAuthService("secret-a", time.Hour)
	token, err := serviceA.Authenticate("tecnico@prefeitura.rio", "painel@2024")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}

	serviceB := NewAuthService("secret-b", time.Hour)
	_, err = serviceB.ValidateToken(token)
	if err == nil {
		t.Fatal("expected error for wrong secret")
	}
}
