package statistics

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/database/testutil"
)

func TestIntegration_StatisticsRepository_Summary(t *testing.T) {
	testutil.SkipIfNoDB(t)

	pool := testutil.Connect(t)
	migrateAndTruncate(t, pool)

	testutil.InsertChild(t, pool, "sum-child-1", "Alice", 8, "Centro", false, nil)
	testutil.InsertChild(t, pool, "sum-child-2", "Bob", 10, "Norte", true, strPtr("admin"))
	testutil.InsertChild(t, pool, "sum-child-3", "Carlos", 7, "Sul", false, nil)

	healthID := testutil.InsertHealth(t, pool, "sum-child-1", false)
	testutil.InsertHealthAlert(t, pool, healthID, "vacinas_atrasadas", "Vacinas Atrasadas")

	eduID := testutil.InsertEducation(t, pool, "sum-child-2", nil, 85)
	testutil.InsertEducationAlert(t, pool, eduID, "frequencia_baixa", "Frequência Baixa")

	saID := testutil.InsertSocialAssistance(t, pool, "sum-child-3", false, false)
	testutil.InsertSocialAssistanceAlert(t, pool, saID, "cadastro_ausente", "Cadastro Ausente")

	repo := NewStatisticsRepository(pool)
	summary, err := repo.Summary(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if summary.TotalChildren != 3 {
		t.Errorf("expected 3 total children, got %d", summary.TotalChildren)
	}
	if summary.Reviewed != 1 {
		t.Errorf("expected 1 reviewed, got %d", summary.Reviewed)
	}
	if summary.PendingReview != 2 {
		t.Errorf("expected 2 pending review, got %d", summary.PendingReview)
	}
	if summary.AlertsByArea["health"] != 1 {
		t.Errorf("expected 1 health alert, got %d", summary.AlertsByArea["health"])
	}
	if summary.AlertsByArea["education"] != 1 {
		t.Errorf("expected 1 education alert, got %d", summary.AlertsByArea["education"])
	}
	if summary.AlertsByArea["social_assistance"] != 1 {
		t.Errorf("expected 1 social assistance alert, got %d", summary.AlertsByArea["social_assistance"])
	}
}

func TestIntegration_StatisticsRepository_Summary_Empty(t *testing.T) {
	testutil.SkipIfNoDB(t)

	pool := testutil.Connect(t)
	migrateAndTruncate(t, pool)

	repo := NewStatisticsRepository(pool)
	summary, err := repo.Summary(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if summary.TotalChildren != 0 {
		t.Errorf("expected 0 total children, got %d", summary.TotalChildren)
	}
	if summary.Reviewed != 0 {
		t.Errorf("expected 0 reviewed, got %d", summary.Reviewed)
	}
	if summary.PendingReview != 0 {
		t.Errorf("expected 0 pending review, got %d", summary.PendingReview)
	}
	for area, count := range summary.AlertsByArea {
		if count != 0 {
			t.Errorf("expected 0 alerts for %s, got %d", area, count)
		}
	}
}

func TestIntegration_StatisticsRepository_StatisticsByNeighborhood(t *testing.T) {
	testutil.SkipIfNoDB(t)

	pool := testutil.Connect(t)
	migrateAndTruncate(t, pool)

	testutil.InsertChild(t, pool, "stat-child-1", "Alice", 8, "Centro", false, nil)
	testutil.InsertChild(t, pool, "stat-child-2", "Bob", 10, "Norte", false, nil)
	testutil.InsertChild(t, pool, "stat-child-3", "Carlos", 7, "Centro", false, nil)

	healthID1 := testutil.InsertHealth(t, pool, "stat-child-1", false)
	testutil.InsertHealthAlert(t, pool, healthID1, "vacinas_atrasadas", "Vacinas Atrasadas")

	healthID2 := testutil.InsertHealth(t, pool, "stat-child-2", false)
	testutil.InsertHealthAlert(t, pool, healthID2, "consulta_atrasada", "Consulta Atrasada")

	eduID3 := testutil.InsertEducation(t, pool, "stat-child-3", nil, 70)
	testutil.InsertEducationAlert(t, pool, eduID3, "matricula_pendente", "Matrícula Pendente")

	repo := NewStatisticsRepository(pool)
	items, err := repo.StatisticsByNeighborhood(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(items) == 0 {
		t.Fatal("expected at least one neighborhood")
	}

	centroFound := false
	for _, item := range items {
		if item.Neighborhood == "Centro" {
			centroFound = true
			if item.Health < 1 {
				t.Errorf("expected at least 1 health alert in Centro, got %d", item.Health)
			}
			if item.Education < 1 {
				t.Errorf("expected at least 1 education alert in Centro, got %d", item.Education)
			}
		}
	}
	if !centroFound {
		t.Error("expected Centro in results")
	}
}

func TestIntegration_StatisticsRepository_StatisticsByNeighborhood_NoAlerts(t *testing.T) {
	testutil.SkipIfNoDB(t)

	pool := testutil.Connect(t)
	migrateAndTruncate(t, pool)

	testutil.InsertChild(t, pool, "no-alert-1", "Alice", 8, "Centro", false, nil)
	testutil.InsertChild(t, pool, "no-alert-2", "Bob", 10, "Norte", false, nil)
	testutil.InsertHealth(t, pool, "no-alert-1", true)
	testutil.InsertHealth(t, pool, "no-alert-2", true)

	repo := NewStatisticsRepository(pool)
	items, err := repo.StatisticsByNeighborhood(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) == 0 {
		t.Fatal("expected neighborhoods even without alerts")
	}
	for _, item := range items {
		if item.Health != 0 || item.Education != 0 || item.SocialAssistance != 0 {
			t.Errorf("expected 0 alerts for %s, got health=%d edu=%d social=%d",
				item.Neighborhood, item.Health, item.Education, item.SocialAssistance)
		}
	}
}

func TestIntegration_StatisticsRepository_StatisticsByNeighborhood_Empty(t *testing.T) {
	testutil.SkipIfNoDB(t)

	pool := testutil.Connect(t)
	migrateAndTruncate(t, pool)

	repo := NewStatisticsRepository(pool)
	items, err := repo.StatisticsByNeighborhood(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if items == nil {
		t.Fatal("expected non-nil slice, got nil")
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}

func strPtr(s string) *string {
	return &s
}

func migrateAndTruncate(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	testutil.Migrate(t, testutil.DSN(t), testutil.MigrationsPath(t))
	testutil.Truncate(t, pool)
}
