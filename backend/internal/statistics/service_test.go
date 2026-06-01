package statistics

import (
	"context"
	"errors"
	"testing"
)

func TestService_GetSummary_Success(t *testing.T) {
	mock := &mockStatisticRepository{
		summaryFunc: func(_ context.Context) (*Summary, error) {
			return &Summary{
				TotalChildren: 100,
				Reviewed:      60,
				PendingReview: 40,
				AlertsByArea:  map[string]int{"health": 30, "education": 20, "social_assistance": 10},
			}, nil
		},
	}
	svc := NewStatisticsService(mock)

	summary, err := svc.GetSummary(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.TotalChildren != 100 {
		t.Errorf("expected 100, got %d", summary.TotalChildren)
	}
	if summary.Reviewed != 60 {
		t.Errorf("expected 60, got %d", summary.Reviewed)
	}
	if summary.PendingReview != 40 {
		t.Errorf("expected 40, got %d", summary.PendingReview)
	}
	if summary.AlertsByArea["health"] != 30 {
		t.Errorf("expected 30 health alerts, got %d", summary.AlertsByArea["health"])
	}
}

func TestService_GetSummary_ZeroValues(t *testing.T) {
	mock := &mockStatisticRepository{
		summaryFunc: func(_ context.Context) (*Summary, error) {
			return &Summary{
				TotalChildren: 0,
				Reviewed:      0,
				PendingReview: 0,
				AlertsByArea:  map[string]int{"health": 0, "education": 0, "social_assistance": 0},
			}, nil
		},
	}
	svc := NewStatisticsService(mock)

	summary, err := svc.GetSummary(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.TotalChildren != 0 {
		t.Errorf("expected 0, got %d", summary.TotalChildren)
	}
}

func TestService_GetSummary_RepositoryError(t *testing.T) {
	mock := &mockStatisticRepository{
		summaryFunc: func(_ context.Context) (*Summary, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewStatisticsService(mock)

	_, err := svc.GetSummary(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestService_GetStatistics_Success(t *testing.T) {
	mock := &mockStatisticRepository{
		statisticsByNeighborhoodFunc: func(_ context.Context) ([]NeighborhoodAlertCount, error) {
			return []NeighborhoodAlertCount{
				{Neighborhood: "Centro", Health: 5, Education: 3, SocialAssistance: 1},
				{Neighborhood: "Norte", Health: 2, Education: 7, SocialAssistance: 0},
			}, nil
		},
	}
	svc := NewStatisticsService(mock)

	result, err := svc.GetStatistics(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Statistics) != 2 {
		t.Errorf("expected 2 items, got %d", len(result.Statistics))
	}
	if result.Statistics[0].Neighborhood != "Centro" {
		t.Errorf("expected Centro, got %s", result.Statistics[0].Neighborhood)
	}
	if result.Statistics[1].Health != 2 {
		t.Errorf("expected 2 health alerts in Norte, got %d", result.Statistics[1].Health)
	}
}

func TestService_GetStatistics_Empty(t *testing.T) {
	mock := &mockStatisticRepository{
		statisticsByNeighborhoodFunc: func(_ context.Context) ([]NeighborhoodAlertCount, error) {
			return []NeighborhoodAlertCount{}, nil
		},
	}
	svc := NewStatisticsService(mock)

	result, err := svc.GetStatistics(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Statistics) != 0 {
		t.Errorf("expected 0 items, got %d", len(result.Statistics))
	}
}

func TestService_GetStatistics_RepositoryError(t *testing.T) {
	mock := &mockStatisticRepository{
		statisticsByNeighborhoodFunc: func(_ context.Context) ([]NeighborhoodAlertCount, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewStatisticsService(mock)

	_, err := svc.GetStatistics(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestService_GetSummary_AlertsByArea_Empty(t *testing.T) {
	mock := &mockStatisticRepository{
		summaryFunc: func(_ context.Context) (*Summary, error) {
			return &Summary{
				TotalChildren: 50,
				Reviewed:      25,
				PendingReview: 25,
				AlertsByArea:  map[string]int{},
			}, nil
		},
	}
	svc := NewStatisticsService(mock)

	summary, err := svc.GetSummary(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(summary.AlertsByArea) != 0 {
		t.Errorf("expected empty alerts map, got %d", len(summary.AlertsByArea))
	}
}
