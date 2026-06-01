package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setupAuthTest() (*AuthHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	service := NewAuthService("test-secret", time.Hour)
	handler := NewAuthHandler(service)
	r := gin.New()
	r.POST("/auth/token", handler.Token)
	r.GET("/auth/session", handler.Session)
	r.DELETE("/auth/session", handler.Logout)
	return handler, r
}

func TestToken_Success(t *testing.T) {
	_, r := setupAuthTest()

	w := httptest.NewRecorder()
	body := `{"email":"tecnico@prefeitura.rio","password":"painel@2024"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/token", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp tokenResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.AccessToken == "" {
		t.Fatal("expected non-empty access_token")
	}
	if resp.TokenType != "Bearer" {
		t.Errorf("expected Bearer, got %s", resp.TokenType)
	}
	if resp.ExpiresIn <= 0 {
		t.Errorf("expected positive expires_in, got %d", resp.ExpiresIn)
	}

	cookies := w.Result().Cookies()
	var found bool
	for _, c := range cookies {
		if c.Name == "auth_token" {
			found = true
			if c.Value == "" {
				t.Error("expected non-empty auth_token cookie")
			}
			if !c.HttpOnly {
				t.Error("expected auth_token cookie to be HttpOnly")
			}
			break
		}
	}
	if !found {
		t.Error("expected auth_token cookie to be set")
	}
}

func TestToken_InvalidCredentials(t *testing.T) {
	_, r := setupAuthTest()

	w := httptest.NewRecorder()
	body := `{"email":"wrong@email.com","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/token", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestToken_InvalidPayload(t *testing.T) {
	_, r := setupAuthTest()

	w := httptest.NewRecorder()
	body := `{"email":""}`
	req := httptest.NewRequest(http.MethodPost, "/auth/token", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestToken_MalformedJSON(t *testing.T) {
	_, r := setupAuthTest()

	w := httptest.NewRecorder()
	body := `not-json`
	req := httptest.NewRequest(http.MethodPost, "/auth/token", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestToken_WrongPassword(t *testing.T) {
	_, r := setupAuthTest()

	tests := []struct {
		name     string
		email    string
		password string
	}{
		{"wrong email", "admin@test.com", "painel@2024"},
		{"wrong password", "tecnico@prefeitura.rio", "wrongpass"},
		{"both wrong", "x@y.com", "z"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			body := `{"email":"` + tt.email + `","password":"` + tt.password + `"}`
			req := httptest.NewRequest(http.MethodPost, "/auth/token", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("expected 401, got %d", w.Code)
			}
		})
	}
}

func TestSession_ValidToken(t *testing.T) {
	service := NewAuthService("test-secret", time.Hour)
	token, err := service.Authenticate("tecnico@prefeitura.rio", "painel@2024")
	if err != nil {
		t.Fatalf("failed to authenticate: %v", err)
	}

	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(service)
	r := gin.New()
	r.GET("/auth/session", handler.Session)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp tokenResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.AccessToken != token {
		t.Errorf("expected same token, got %s", resp.AccessToken)
	}
	if resp.TokenType != "Bearer" {
		t.Errorf("expected Bearer, got %s", resp.TokenType)
	}
	if resp.ExpiresIn <= 0 {
		t.Errorf("expected positive expires_in, got %d", resp.ExpiresIn)
	}
}

func TestSession_NoCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := NewAuthService("test-secret", time.Hour)
	handler := NewAuthHandler(service)
	r := gin.New()
	r.GET("/auth/session", handler.Session)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestSession_InvalidCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := NewAuthService("test-secret", time.Hour)
	handler := NewAuthHandler(service)
	r := gin.New()
	r.GET("/auth/session", handler.Session)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: "invalid-token"})
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestSession_ExpiredToken(t *testing.T) {
	service := NewAuthService("test-secret", -time.Hour)
	token, err := GenerateToken("tecnico@prefeitura.rio", "test-secret", -time.Hour)
	if err != nil {
		t.Fatalf("failed to generate expired token: %v", err)
	}

	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(service)
	r := gin.New()
	r.GET("/auth/session", handler.Session)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLogout(t *testing.T) {
	_, r := setupAuthTest()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/auth/session", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}

	cookies := w.Result().Cookies()
	var found bool
	for _, c := range cookies {
		if c.Name == "auth_token" {
			found = true
			if c.Value != "" {
				t.Error("expected empty auth_token cookie value after logout")
			}
			if c.MaxAge != -1 {
				t.Errorf("expected MaxAge -1, got %d", c.MaxAge)
			}
			break
		}
	}
	if !found {
		t.Error("expected auth_token cookie in response")
	}
}

func TestToken_MissingContentType(t *testing.T) {
	_, r := setupAuthTest()

	w := httptest.NewRecorder()
	body := `not-json`
	req := httptest.NewRequest(http.MethodPost, "/auth/token", strings.NewReader(body))
	req.Header.Set("Content-Type", "text/plain")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestToken_EmptyBody(t *testing.T) {
	_, r := setupAuthTest()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/token", nil)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
