package statistics

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupStatsHandler() (*StatisticsHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mock := &mockStatisticRepository{
		summaryFunc: func(ctx context.Context) (*Summary, error) {
			return &Summary{
				TotalChildren: 100,
				Reviewed:      60,
				PendingReview: 40,
				AlertsByArea:  map[string]int{"health": 30, "education": 20, "social_assistance": 10},
			}, nil
		},
		statisticsByNeighborhoodFunc: func(ctx context.Context) ([]NeighborhoodAlertCount, error) {
			return []NeighborhoodAlertCount{
				{Neighborhood: "Centro", Health: 5, Education: 3, SocialAssistance: 1},
			}, nil
		},
	}
	svc := NewStatisticsService(mock)
	handler := NewStatisticsHandler(svc)
	r := gin.New()
	r.GET("/summary", handler.GetSummary)
	r.GET("/statistics", handler.GetStatistics)
	return handler, r
}

func TestHandler_GetSummary_Success(t *testing.T) {
	_, r := setupStatsHandler()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/summary", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var summary Summary
	if err := json.Unmarshal(w.Body.Bytes(), &summary); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if summary.TotalChildren != 100 {
		t.Errorf("expected 100, got %d", summary.TotalChildren)
	}
	if summary.AlertsByArea["health"] != 30 {
		t.Errorf("expected 30 health alerts, got %d", summary.AlertsByArea["health"])
	}
}

func TestHandler_GetSummary_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockStatisticRepository{
		summaryFunc: func(ctx context.Context) (*Summary, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewStatisticsService(mock)
	handler := NewStatisticsHandler(svc)
	r := gin.New()
	r.GET("/summary", handler.GetSummary)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/summary", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestHandler_GetStatistics_Success(t *testing.T) {
	_, r := setupStatsHandler()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/statistics", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp StatisticsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(resp.Statistics) != 1 {
		t.Errorf("expected 1 item, got %d", len(resp.Statistics))
	}
}

func TestHandler_GetStatistics_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockStatisticRepository{
		statisticsByNeighborhoodFunc: func(ctx context.Context) ([]NeighborhoodAlertCount, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewStatisticsService(mock)
	handler := NewStatisticsHandler(svc)
	r := gin.New()
	r.GET("/statistics", handler.GetStatistics)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/statistics", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestHandler_GetSummary_ZeroValues(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockStatisticRepository{
		summaryFunc: func(ctx context.Context) (*Summary, error) {
			return &Summary{
				TotalChildren: 0,
				Reviewed:      0,
				PendingReview: 0,
				AlertsByArea:  map[string]int{"health": 0, "education": 0, "social_assistance": 0},
			}, nil
		},
	}
	svc := NewStatisticsService(mock)
	handler := NewStatisticsHandler(svc)
	r := gin.New()
	r.GET("/summary", handler.GetSummary)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/summary", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var summary Summary
	if err := json.Unmarshal(w.Body.Bytes(), &summary); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if summary.TotalChildren != 0 {
		t.Errorf("expected 0, got %d", summary.TotalChildren)
	}
}
