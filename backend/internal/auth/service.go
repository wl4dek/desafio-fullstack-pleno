package auth

import (
	"errors"
	"time"
)

type AuthService struct {
	secret string
	expiry time.Duration
}

func NewAuthService(secret string, expiry time.Duration) *AuthService {
	return &AuthService{secret: secret, expiry: expiry}
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (s *AuthService) Authenticate(email, password string) (string, error) {
	if email != "tecnico@prefeitura.rio" || password != "painel@2024" {
		return "", ErrInvalidCredentials
	}

	token, err := GenerateToken(email, s.secret, s.expiry)
	if err != nil {
		return "", err
	}

	return token, nil
}
