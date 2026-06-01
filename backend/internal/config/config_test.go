package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	os.Clearenv()

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("expected port 8080, got %s", cfg.Port)
	}
	if cfg.DatabaseURL != "postgres://postgres:postgres@localhost:5432/full-stack?sslmode=disable" {
		t.Errorf("unexpected default database URL: %s", cfg.DatabaseURL)
	}
	if cfg.JWTSecret != "super-secret" {
		t.Errorf("expected super-secret, got %s", cfg.JWTSecret)
	}
	if len(cfg.AllowOrigins) != 1 || cfg.AllowOrigins[0] != "http://localhost:3000" {
		t.Errorf("unexpected AllowOrigins: %v", cfg.AllowOrigins)
	}
	if cfg.BasePath != "." {
		t.Errorf("expected base path '.', got %s", cfg.BasePath)
	}
}

func TestLoad_CustomPort(t *testing.T) {
	os.Clearenv()
	os.Setenv("PORT", "9090")
	defer os.Unsetenv("PORT")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("expected 9090, got %s", cfg.Port)
	}
}

func TestLoad_CustomDatabaseURL(t *testing.T) {
	os.Clearenv()
	dsn := "postgres://user:pass@host:5432/db?sslmode=disable"
	os.Setenv("DATABASE_URL", dsn)
	defer os.Unsetenv("DATABASE_URL")

	cfg := Load()

	if cfg.DatabaseURL != dsn {
		t.Errorf("expected custom DSN, got %s", cfg.DatabaseURL)
	}
}

func TestLoad_CustomJWTSecret(t *testing.T) {
	os.Clearenv()
	os.Setenv("JWT_SECRET", "my-custom-secret")
	defer os.Unsetenv("JWT_SECRET")

	cfg := Load()

	if cfg.JWTSecret != "my-custom-secret" {
		t.Errorf("expected my-custom-secret, got %s", cfg.JWTSecret)
	}
}

func TestLoad_CustomAllowOrigins(t *testing.T) {
	os.Clearenv()
	os.Setenv("ALLOW_ORIGINS", "http://localhost:3000,https://app.example.com")
	defer os.Unsetenv("ALLOW_ORIGINS")

	cfg := Load()

	if len(cfg.AllowOrigins) != 2 {
		t.Fatalf("expected 2 origins, got %d", len(cfg.AllowOrigins))
	}
	if cfg.AllowOrigins[0] != "http://localhost:3000" {
		t.Errorf("expected http://localhost:3000, got %s", cfg.AllowOrigins[0])
	}
	if cfg.AllowOrigins[1] != "https://app.example.com" {
		t.Errorf("expected https://app.example.com, got %s", cfg.AllowOrigins[1])
	}
}

func TestLoad_CustomBasePath(t *testing.T) {
	os.Clearenv()
	os.Setenv("BASE_PATH", "/app")
	defer os.Unsetenv("BASE_PATH")

	cfg := Load()

	if cfg.BasePath != "/app" {
		t.Errorf("expected /app, got %s", cfg.BasePath)
	}
}

func TestLoad_AllEnvVars(t *testing.T) {
	os.Clearenv()
	os.Setenv("PORT", "3000")
	os.Setenv("DATABASE_URL", "postgres://custom:custom@localhost:5432/mydb?sslmode=disable")
	os.Setenv("JWT_SECRET", "custom-jwt")
	os.Setenv("ALLOW_ORIGINS", "*")
	os.Setenv("BASE_PATH", "/custom/path")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("ALLOW_ORIGINS")
		os.Unsetenv("BASE_PATH")
	}()

	cfg := Load()

	if cfg.Port != "3000" {
		t.Errorf("expected 3000, got %s", cfg.Port)
	}
	if cfg.JWTSecret != "custom-jwt" {
		t.Errorf("expected custom-jwt, got %s", cfg.JWTSecret)
	}
	if len(cfg.AllowOrigins) != 1 || cfg.AllowOrigins[0] != "*" {
		t.Errorf("expected *, got %v", cfg.AllowOrigins)
	}
	if cfg.BasePath != "/custom/path" {
		t.Errorf("expected /custom/path, got %s", cfg.BasePath)
	}
}
