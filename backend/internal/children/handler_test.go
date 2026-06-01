package children

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupHandlerTest() (*ChildHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mock := &mockChildRepository{
		listFunc: func(_ context.Context, _ Filters) ([]Child, error) {
			return []Child{
				{ID: "1", Name: "Alice", Age: 8, Neighborhood: "Centro", Reviewed: false},
			}, nil
		},
		countFilteredFunc: func(_ context.Context, _ Filters) (int, error) {
			return 1, nil
		},
		findByIDFunc: func(_ context.Context, id string) (*ChildById, error) {
			return &ChildById{
				Child: Child{ID: id, Name: "Bob", Age: 10, Neighborhood: "Norte"},
			}, nil
		},
		markReviewedFunc: func(_ context.Context, id string, _ string) error {
			return nil
		},
		listNeighborhoodFunc: func(_ context.Context) ([]string, error) {
			return []string{"Centro", "Norte"}, nil
		},
		listAlertsByChildIDFunc: func(_ context.Context, _ string) ([]Alerts, error) {
			return []Alerts{}, nil
		},
	}
	svc := NewChildService(mock)
	handler := NewChildHandler(svc)
	r := gin.New()
	r.GET("/children", handler.List)
	r.GET("/children/neighborhood", handler.ListNeighborhood)
	r.GET("/children/:id", handler.GetByID)
	r.PATCH("/children/:id/review", handler.MarkReviewed)
	return handler, r
}

func TestList_DefaultParams(t *testing.T) {
	_, r := setupHandlerTest()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/children", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp PaginatedResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Errorf("expected 1 child, got %d", len(resp.Data))
	}
	if resp.Pagination.Total != 1 {
		t.Errorf("expected total 1, got %d", resp.Pagination.Total)
	}
}

func TestList_WithQueryParams(t *testing.T) {
	_, r := setupHandlerTest()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/children?neighborhood=Centro&page=2&per_page=20", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetByID_Found(t *testing.T) {
	_, r := setupHandlerTest()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/children/child-1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp ChildByIdResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if resp.ID != "child-1" {
		t.Errorf("expected child-1, got %s", resp.ID)
	}
}

func TestMarkReviewed_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockChildRepository{
		markReviewedFunc: func(_ context.Context, id string, _ string) error {
			return nil
		},
	}
	svc := NewChildService(mock)
	handler := NewChildHandler(svc)
	r := gin.New()
	r.PATCH("/children/:id/review", handler.MarkReviewed)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/children/child-1/review", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestListNeighborhood(t *testing.T) {
	_, r := setupHandlerTest()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/children/neighborhood", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var neighborhoods []string
	if err := json.Unmarshal(w.Body.Bytes(), &neighborhoods); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(neighborhoods) != 2 {
		t.Errorf("expected 2 neighborhoods, got %d", len(neighborhoods))
	}
}

func TestGetByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockChildRepository{
		findByIDFunc: func(_ context.Context, id string) (*ChildById, error) {
			return nil, nil
		},
	}
	svc := NewChildService(mock)
	handler := NewChildHandler(svc)
	r := gin.New()
	r.GET("/children/:id", handler.GetByID)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/children/nonexistent", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetByID_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockChildRepository{
		findByIDFunc: func(_ context.Context, id string) (*ChildById, error) {
			return nil, errInternal
		},
		listAlertsByChildIDFunc: func(_ context.Context, _ string) ([]Alerts, error) {
			return nil, nil
		},
	}
	svc := NewChildService(mock)
	handler := NewChildHandler(svc)
	r := gin.New()
	r.GET("/children/:id", handler.GetByID)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/children/child-1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestList_InvalidPerPage(t *testing.T) {
	_, r := setupHandlerTest()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/children?per_page=5", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestList_InvalidPage(t *testing.T) {
	_, r := setupHandlerTest()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/children?page=-1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestList_WithAlertFilter(t *testing.T) {
	_, r := setupHandlerTest()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/children?alert=vacinas_atrasadas", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestList_WithReviewedFilter(t *testing.T) {
	_, r := setupHandlerTest()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/children?reviewed=true", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestList_WithHasAlertFilter(t *testing.T) {
	_, r := setupHandlerTest()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/children?has_alert=true", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestMarkReviewed_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockChildRepository{
		markReviewedFunc: func(_ context.Context, id string, _ string) error {
			return ErrChildNotFound
		},
	}
	svc := NewChildService(mock)
	handler := NewChildHandler(svc)
	r := gin.New()
	r.PATCH("/children/:id/review", handler.MarkReviewed)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/children/nonexistent/review", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

var errInternal = &internalError{}

type internalError struct{}

func (e *internalError) Error() string { return "internal" }
