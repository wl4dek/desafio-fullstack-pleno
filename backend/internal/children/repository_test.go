package children_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/children"
	"backend/internal/database/testutil"
)

func TestIntegration_ChildRepository_List(t *testing.T) {
	testutil.SkipIfNoDB(t)

	pool := testutil.Connect(t)
	migrateAndTruncate(t, pool)

	testutil.InsertChild(t, pool, "hood-test", "Test", 5, "Centro", false, nil)

	var count int
	if err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM children WHERE id = 'hood-test'").Scan(&count); err != nil {
		t.Fatalf("failed to verify insert: %v", err)
	}
	t.Logf("child count after insert: %d", count)

	repo := children.NewChildRepository(pool)
	result, err := repo.List(context.Background(), children.Filters{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Fatal("expected at least one child")
	}
}

func TestIntegration_ChildRepository_CountFiltered(t *testing.T) {
	testutil.SkipIfNoDB(t)

	pool := testutil.Connect(t)
	migrateAndTruncate(t, pool)

	repo := children.NewChildRepository(pool)
	count, err := repo.CountFiltered(context.Background(), children.Filters{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count < 0 {
		t.Errorf("expected non-negative count, got %d", count)
	}
}

func TestIntegration_ChildRepository_FindByID(t *testing.T) {
	testutil.SkipIfNoDB(t)

	pool := testutil.Connect(t)
	migrateAndTruncate(t, pool)

	testutil.InsertChild(t, pool, "test-id-1", "Test Child", 8, "Centro", false, nil)

	repo := children.NewChildRepository(pool)
	child, err := repo.FindByID(context.Background(), "test-id-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if child == nil {
		t.Fatal("expected child to be found")
	}
	if child.Name != "Test Child" {
		t.Errorf("expected Test Child, got %s", child.Name)
	}
	if child.Neighborhood != "Centro" {
		t.Errorf("expected Centro, got %s", child.Neighborhood)
	}
}

func TestIntegration_ChildRepository_FindByID_NotFound(t *testing.T) {
	testutil.SkipIfNoDB(t)

	pool := testutil.Connect(t)
	migrateAndTruncate(t, pool)

	repo := children.NewChildRepository(pool)
	child, err := repo.FindByID(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if child != nil {
		t.Fatal("expected nil for non-existent child")
	}
}

func TestIntegration_ChildRepository_MarkReviewed(t *testing.T) {
	testutil.SkipIfNoDB(t)

	pool := testutil.Connect(t)
	migrateAndTruncate(t, pool)

	testutil.InsertChild(t, pool, "review-test-1", "Review Child", 10, "Norte", false, nil)

	repo := children.NewChildRepository(pool)
	err := repo.MarkReviewed(context.Background(), "review-test-1", "admin@test.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	child, _ := repo.FindByID(context.Background(), "review-test-1")
	if child == nil {
		t.Fatal("expected child to exist")
	}
	if !child.Reviewed {
		t.Error("expected child to be marked as reviewed")
	}
	if child.ReviewedBy == nil || *child.ReviewedBy != "admin@test.com" {
		t.Errorf("expected reviewed_by admin@test.com, got %v", child.ReviewedBy)
	}
}

func TestIntegration_ChildRepository_MarkReviewed_NotFound(t *testing.T) {
	testutil.SkipIfNoDB(t)

	pool := testutil.Connect(t)
	migrateAndTruncate(t, pool)

	repo := children.NewChildRepository(pool)
	err := repo.MarkReviewed(context.Background(), "nonexistent", "admin")
	if err != children.ErrChildNotFound {
		t.Errorf("expected ErrChildNotFound, got %v", err)
	}
}

func TestIntegration_ChildRepository_MarkReviewed_AlreadyReviewed(t *testing.T) {
	testutil.SkipIfNoDB(t)

	pool := testutil.Connect(t)
	migrateAndTruncate(t, pool)

	admin := "admin@test.com"
	testutil.InsertChild(t, pool, "already-reviewed", "Done Child", 7, "Sul", true, &admin)

	repo := children.NewChildRepository(pool)
	err := repo.MarkReviewed(context.Background(), "already-reviewed", "other@test.com")
	if err != children.ErrChildNotFound {
		t.Errorf("expected ErrChildNotFound for already reviewed child, got %v", err)
	}
}

func TestIntegration_ChildRepository_ListNeighborhood(t *testing.T) {
	testutil.SkipIfNoDB(t)

	pool := testutil.Connect(t)
	migrateAndTruncate(t, pool)

	testutil.InsertChild(t, pool, "hood-test", "Test", 5, "Centro", false, nil)

	repo := children.NewChildRepository(pool)
	neighborhoods, err := repo.ListNeighborhood(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("neighborhoods: %v", neighborhoods)
	if len(neighborhoods) == 0 {
		t.Error("expected at least one neighborhood")
	}
}

func TestIntegration_ChildRepository_ListAlertsByChildID(t *testing.T) {
	testutil.SkipIfNoDB(t)

	pool := testutil.Connect(t)
	migrateAndTruncate(t, pool)

	testutil.InsertChild(t, pool, "alert-child", "Alert Child", 6, "Leste", false, nil)
	healthID := testutil.InsertHealth(t, pool, "alert-child", false)
	testutil.InsertHealthAlert(t, pool, healthID, "vacinas_atrasadas", "Vacinas Atrasadas")

	repo := children.NewChildRepository(pool)
	alerts, err := repo.ListAlertsByChildID(context.Background(), "alert-child")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(alerts) == 0 {
		t.Fatal("expected at least one alert")
	}
	found := false
	for _, a := range alerts {
		if a.Code == "vacinas_atrasadas" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected vacinas_atrasadas alert")
	}
}

func TestIntegration_ChildRepository_ListAlertsByChildID_NoAlerts(t *testing.T) {
	testutil.SkipIfNoDB(t)

	pool := testutil.Connect(t)
	migrateAndTruncate(t, pool)

	testutil.InsertChild(t, pool, "no-alert-child", "No Alert", 5, "Oeste", false, nil)
	testutil.InsertHealth(t, pool, "no-alert-child", true)

	repo := children.NewChildRepository(pool)
	alerts, err := repo.ListAlertsByChildID(context.Background(), "no-alert-child")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(alerts) != 0 {
		t.Errorf("expected 0 alerts, got %d", len(alerts))
	}
}

func TestIntegration_ChildRepository_Count(t *testing.T) {
	testutil.SkipIfNoDB(t)

	pool := testutil.Connect(t)
	migrateAndTruncate(t, pool)

	repo := children.NewChildRepository(pool)
	count, err := repo.Count(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count < 0 {
		t.Errorf("expected non-negative count, got %d", count)
	}
}

func TestIntegration_ChildRepository_List_WithFilters(t *testing.T) {
	testutil.SkipIfNoDB(t)

	pool := testutil.Connect(t)
	migrateAndTruncate(t, pool)

	testutil.InsertChild(t, pool, "filter-a", "Alice", 8, "Centro", false, nil)
	testutil.InsertChild(t, pool, "filter-b", "Bob", 10, "Norte", true, strRef("admin"))

	repo := children.NewChildRepository(pool)
	result, err := repo.List(context.Background(), children.Filters{
		Neighborhood: "Centro",
		Page:         1,
		PerPage:      10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 child in Centro, got %d", len(result))
	}
}

func TestIntegration_ChildRepository_List_Pagination(t *testing.T) {
	testutil.SkipIfNoDB(t)

	pool := testutil.Connect(t)
	migrateAndTruncate(t, pool)

	for i := 0; i < 15; i++ {
		testutil.InsertChild(t, pool, "page-"+string(rune('a'+i)), "Child", 5, "Centro", false, nil)
	}

	repo := children.NewChildRepository(pool)
	page1, err := repo.List(context.Background(), children.Filters{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(page1) > 10 {
		t.Errorf("expected at most 10 children on page 1, got %d", len(page1))
	}
}

func strRef(s string) *string {
	return &s
}

func migrateAndTruncate(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	testutil.Migrate(t, testutil.DSN(t), testutil.MigrationsPath(t))
	testutil.Truncate(t, pool)
}
