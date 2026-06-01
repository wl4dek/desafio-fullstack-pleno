package children

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
)

func TestChildService_List_Success(t *testing.T) {
	mock := &mockChildRepository{
		listFunc: func(_ context.Context, f Filters) ([]Child, error) {
			return []Child{
				{ID: "1", Name: "Alice", Age: 8, Neighborhood: "Centro", Reviewed: false},
			}, nil
		},
		countFilteredFunc: func(_ context.Context, f Filters) (int, error) {
			return 1, nil
		},
	}
	svc := NewChildService(mock)

	result, err := svc.List(context.Background(), Filters{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result.Data) != 1 {
		t.Errorf("expected 1 child, got %d", len(result.Data))
	}
	if result.Pagination.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Pagination.Total)
	}
	if result.Pagination.TotalPages != 1 {
		t.Errorf("expected total_pages 1, got %d", result.Pagination.TotalPages)
	}
}

func TestChildService_List_Pagination(t *testing.T) {
	mock := &mockChildRepository{
		listFunc: func(_ context.Context, f Filters) ([]Child, error) {
			return []Child{{ID: "1", Name: "A"}}, nil
		},
		countFilteredFunc: func(_ context.Context, f Filters) (int, error) {
			return 25, nil
		},
	}
	svc := NewChildService(mock)

	result, err := svc.List(context.Background(), Filters{Page: 3, PerPage: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Pagination.TotalPages != 3 {
		t.Errorf("expected 3 total_pages, got %d", result.Pagination.TotalPages)
	}
	if result.Pagination.Page != 3 {
		t.Errorf("expected page 3, got %d", result.Pagination.Page)
	}
}

func TestChildService_List_Empty(t *testing.T) {
	mock := &mockChildRepository{
		listFunc: func(_ context.Context, f Filters) ([]Child, error) {
			return []Child{}, nil
		},
		countFilteredFunc: func(_ context.Context, f Filters) (int, error) {
			return 0, nil
		},
	}
	svc := NewChildService(mock)

	result, err := svc.List(context.Background(), Filters{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 0 {
		t.Errorf("expected 0 children, got %d", len(result.Data))
	}
	if result.Pagination.Total != 0 {
		t.Errorf("expected total 0, got %d", result.Pagination.Total)
	}
}

func TestChildService_List_RepoError(t *testing.T) {
	mock := &mockChildRepository{
		listFunc: func(_ context.Context, f Filters) ([]Child, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewChildService(mock)

	_, err := svc.List(context.Background(), Filters{Page: 1, PerPage: 10})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestChildService_GetByID_Found(t *testing.T) {
	now := time.Now()
	mock := &mockChildRepository{
		findByIDFunc: func(_ context.Context, id string) (*ChildById, error) {
			return &ChildById{
				Child: Child{
					ID: id, Name: "Bob", Age: 10,
					Neighborhood: "Norte", Reviewed: false,
					CreatedAt: now,
				},
				Health:           Health{VaccinationsUpToDate: true, Alerts: nil},
				Education:        Education{SchoolName: strPtr("Escola X"), FrequenciaPercent: 90, Alerts: nil},
				SocialAssistance: SocialAssistance{CadUnico: true, ActiveBenefit: false, Alerts: nil},
			}, nil
		},
		listAlertsByChildIDFunc: func(_ context.Context, id string) ([]Alerts, error) {
			return []Alerts{}, nil
		},
	}
	svc := NewChildService(mock)

	child, err := svc.GetByID(context.Background(), "child-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if child == nil {
		t.Fatal("expected non-nil child")
	}
	if child.ID != "child-1" {
		t.Errorf("expected child-1, got %s", child.ID)
	}
	if child.Health.VaccinationsUpToDate != true {
		t.Error("expected health data")
	}
}

func TestChildService_GetByID_NotFound(t *testing.T) {
	mock := &mockChildRepository{
		findByIDFunc: func(_ context.Context, id string) (*ChildById, error) {
			return nil, nil
		},
	}
	svc := NewChildService(mock)

	child, err := svc.GetByID(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if child != nil {
		t.Fatal("expected nil child for not found")
	}
}

func TestChildService_GetByID_WithAlerts(t *testing.T) {
	now := time.Now()
	mock := &mockChildRepository{
		findByIDFunc: func(_ context.Context, id string) (*ChildById, error) {
			return &ChildById{
				Child: Child{ID: id, Name: "Carlos", Age: 7, Neighborhood: "Sul", Reviewed: false, CreatedAt: now},
				Health:           Health{VaccinationsUpToDate: false, Alerts: nil},
				Education:        Education{SchoolName: strPtr("Escola Y"), FrequenciaPercent: 80, Alerts: nil},
				SocialAssistance: SocialAssistance{CadUnico: false, ActiveBenefit: false, Alerts: nil},
			}, nil
		},
		listAlertsByChildIDFunc: func(_ context.Context, id string) ([]Alerts, error) {
			return []Alerts{
				{Category: "health", Code: "vacinas_atrasadas", Message: sql.NullString{String: "Vacinas Atrasadas", Valid: true}},
				{Category: "education", Code: "frequencia_baixa", Message: sql.NullString{String: "Frequência Baixa", Valid: true}},
			}, nil
		},
	}
	svc := NewChildService(mock)

	child, err := svc.GetByID(context.Background(), "child-2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if child == nil {
		t.Fatal("expected non-nil child")
	}
	if len(child.Health.Alerts) != 1 || child.Health.Alerts[0] != "Vacinas Atrasadas" {
		t.Errorf("expected health alert, got %v", child.Health.Alerts)
	}
	if len(child.Education.Alerts) != 1 || child.Education.Alerts[0] != "Frequência Baixa" {
		t.Errorf("expected education alert, got %v", child.Education.Alerts)
	}
	if len(child.SocialAssistance.Alerts) != 0 {
		t.Errorf("expected 0 social assistance alerts, got %d", len(child.SocialAssistance.Alerts))
	}
}

func TestChildService_GetByID_RepoError(t *testing.T) {
	mock := &mockChildRepository{
		findByIDFunc: func(_ context.Context, id string) (*ChildById, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewChildService(mock)

	_, err := svc.GetByID(context.Background(), "child-1")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestChildService_MarkReviewed_Success(t *testing.T) {
	mock := &mockChildRepository{
		markReviewedFunc: func(_ context.Context, id string, reviewedBy string) error {
			return nil
		},
	}
	svc := NewChildService(mock)

	err := svc.MarkReviewed(context.Background(), "child-1", "admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestChildService_MarkReviewed_NotFound(t *testing.T) {
	mock := &mockChildRepository{
		markReviewedFunc: func(_ context.Context, id string, reviewedBy string) error {
			return ErrChildNotFound
		},
	}
	svc := NewChildService(mock)

	err := svc.MarkReviewed(context.Background(), "nonexistent", "admin")
	if err != ErrChildNotFound {
		t.Errorf("expected ErrChildNotFound, got %v", err)
	}
}

func TestChildService_ListNeighborhood_Success(t *testing.T) {
	mock := &mockChildRepository{
		listNeighborhoodFunc: func(_ context.Context) ([]string, error) {
			return []string{"Centro", "Norte", "Sul"}, nil
		},
	}
	svc := NewChildService(mock)

	neighborhoods, err := svc.ListNeighborhood(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(neighborhoods) != 3 {
		t.Errorf("expected 3 neighborhoods, got %d", len(neighborhoods))
	}
}

func TestChildService_ListNeighborhood_Empty(t *testing.T) {
	mock := &mockChildRepository{
		listNeighborhoodFunc: func(_ context.Context) ([]string, error) {
			return []string{}, nil
		},
	}
	svc := NewChildService(mock)

	neighborhoods, err := svc.ListNeighborhood(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(neighborhoods) != 0 {
		t.Errorf("expected 0 neighborhoods, got %d", len(neighborhoods))
	}
}

func TestChildService_MarkReviewed_RepoError(t *testing.T) {
	mock := &mockChildRepository{
		markReviewedFunc: func(_ context.Context, id string, reviewedBy string) error {
			return errors.New("db error")
		},
	}
	svc := NewChildService(mock)

	err := svc.MarkReviewed(context.Background(), "child-1", "admin")
	if err == nil {
		t.Fatal("expected error")
	}
}

func strPtr(s string) *string {
	return &s
}

func TestChildService_GetByID_AlertCategories(t *testing.T) {
	now := time.Now()
	mock := &mockChildRepository{
		findByIDFunc: func(_ context.Context, id string) (*ChildById, error) {
			return &ChildById{
				Child: Child{ID: id, Name: "D", Age: 5, Neighborhood: "Leste",
				AlertCategories: []string{"health", "social_assistance"}, Reviewed: false, CreatedAt: now},
				Health:           Health{VaccinationsUpToDate: false},
				Education:        Education{FrequenciaPercent: 70},
				SocialAssistance: SocialAssistance{CadUnico: false, ActiveBenefit: false},
			}, nil
		},
		listAlertsByChildIDFunc: func(_ context.Context, id string) ([]Alerts, error) {
			return []Alerts{
				{Category: "health", Code: "consulta_atrasada", Message: sql.NullString{String: "Consulta Atrasada", Valid: true}},
				{Category: "social_assistance", Code: "cadastro_ausente", Message: sql.NullString{String: "Cadastro Ausente", Valid: true}},
			}, nil
		},
	}
	svc := NewChildService(mock)

	child, err := svc.GetByID(context.Background(), "child-3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if child == nil {
		t.Fatal("expected non-nil child")
	}

	foundHealth, foundSocial := false, false
	for _, ac := range child.AlertCategories {
		if ac == "health" {
			foundHealth = true
		}
		if ac == "social_assistance" {
			foundSocial = true
		}
	}
	if !foundHealth {
		t.Error("expected health alert category")
	}
	if !foundSocial {
		t.Error("expected social_assistance alert category")
	}
}
