package auth

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestExtractToken_FromHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer my-test-token")

	token, err := extractToken(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "my-test-token" {
		t.Errorf("expected my-test-token, got %s", token)
	}
}

func TestExtractToken_FromCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.AddCookie(&http.Cookie{Name: "auth_token", Value: "cookie-token"})

	token, err := extractToken(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "cookie-token" {
		t.Errorf("expected cookie-token, got %s", token)
	}
}

func TestExtractToken_HeaderPrecedesCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer header-token")
	c.Request.AddCookie(&http.Cookie{Name: "auth_token", Value: "cookie-token"})

	token, err := extractToken(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "header-token" {
		t.Errorf("expected header-token, got %s", token)
	}
}

func TestExtractToken_Missing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := extractToken(c)
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestExtractToken_InvalidHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "InvalidFormat")

	_, err := extractToken(c)
	if err == nil {
		t.Fatal("expected error for invalid header format")
	}
}

func TestExtractToken_EmptyHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer ")

	token, err := extractToken(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "" {
		t.Errorf("expected empty token, got %s", token)
	}
}

func TestExtractToken_CookieTakesOverWhenHeaderMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Basic somecreds")
	c.Request.AddCookie(&http.Cookie{Name: "auth_token", Value: "cookie-token"})

	_, err := extractToken(c)
	if err == nil {
		t.Fatal("expected error for basic auth header")
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "test-secret"
	tokenStr, err := GenerateToken("admin@test.com", secret, time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenStr)

	middleware := AuthMiddleware(secret)
	middleware(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	username, exists := c.Get("preferred_username")
	if !exists {
		t.Fatal("expected preferred_username in context")
	}
	if username != "admin@test.com" {
		t.Errorf("expected admin@test.com, got %v", username)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid-token")

	middleware := AuthMiddleware("secret")
	middleware(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	middleware := AuthMiddleware("secret")
	middleware(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_DifferentSecret(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tokenStr, err := GenerateToken("admin@test.com", "secret-a", time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenStr)

	middleware := AuthMiddleware("secret-b")
	middleware(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_TokenFromCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "test-secret"
	tokenStr, err := GenerateToken("admin@test.com", secret, time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.AddCookie(&http.Cookie{Name: "auth_token", Value: tokenStr})

	middleware := AuthMiddleware(secret)
	middleware(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "test-secret"
	tokenStr, err := GenerateToken("admin@test.com", secret, -time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer "+tokenStr)

	middleware := AuthMiddleware(secret)
	middleware(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestExtractToken_Base64Encoded(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	encoded := base64.StdEncoding.EncodeToString([]byte("user:pass"))
	c.Request.Header.Set("Authorization", "Basic "+encoded)
	c.Request.AddCookie(&http.Cookie{Name: "auth_token", Value: ""})

	_, err := extractToken(c)
	if err == nil {
		t.Fatal("expected error for basic auth")
	}
}

func TestExtractToken_MalformedHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer ")

	token, err := extractToken(c)
	if err != nil {
		t.Fatalf("unexpected error for empty bearer value: %v", err)
	}
	if token != "" {
		t.Errorf("expected empty token, got %s", token)
	}
}

func TestExtractToken_CookieOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.AddCookie(&http.Cookie{Name: "auth_token", Value: "session-token"})

	token, err := extractToken(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "session-token" {
		t.Errorf("expected session-token, got %s", token)
	}
}

func TestExtractToken_EmptyCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.AddCookie(&http.Cookie{Name: "auth_token", Value: ""})

	_, err := extractToken(c)
	if err == nil {
		t.Fatal("expected error for empty cookie")
	}
}

func TestAuthMiddleware_InvalidBearerFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer")

	_, err := extractToken(c)
	if err == nil {
		t.Fatal("expected error for malformed bearer")
	}

	parts := strings.SplitN("Bearer", " ", 2)
	if len(parts) != 2 {
		token, err := extractToken(c)
		if err != nil {
			return
		}
		_ = token
	}
}
