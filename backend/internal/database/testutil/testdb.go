package testutil

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/database"
)

func Connect(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := DSN(t)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func Migrate(t *testing.T, dsn, migrationsPath string) {
	t.Helper()
	if err := database.RunMigrations(dsn, migrationsPath); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}
}

func Truncate(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `
		TRUNCATE TABLE
			alert_health, alert_education, alert_social_assistance,
			health, education, social_assistance, children CASCADE`)
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
}

func Seed(t *testing.T, pool *pgxpool.Pool, seedPath string) {
	t.Helper()
	if err := database.LoadSeed(context.Background(), pool, seedPath); err != nil {
		t.Fatalf("failed to load seed data: %v", err)
	}
}

func MigrationsPath(t *testing.T) string {
	t.Helper()
	return resolvePath(t, "migrations")
}

func SeedPath(t *testing.T) string {
	t.Helper()
	return resolvePath(t, "data/seed.json")
}

func DSN(t *testing.T) string {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/full-stack?sslmode=disable"
	}
	return dsn
}

func SkipIfNoDB(t *testing.T) {
	t.Helper()
	pool, err := pgxpool.New(context.Background(), DSN(t))
	if err != nil {
		t.Skipf("cannot connect to test database: %v", err)
	}
	pool.Close()
}

func resolvePath(t *testing.T, target string) string {
	t.Helper()

	base := os.Getenv("BASE_PATH")
	if base != "" {
		abs, err := filepath.Abs(filepath.Join(base, target))
		if err == nil {
			return abs
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	dir := wd
	for {
		candidate := filepath.Join(dir, target)
		if _, err := os.Stat(candidate); err == nil {
			abs, _ := filepath.Abs(candidate)
			return abs
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	t.Fatalf("could not find %s from %s", target, wd)
	return ""
}

func InsertChild(t *testing.T, pool *pgxpool.Pool, id, name string, age int, neighborhood string, reviewed bool, reviewedBy *string) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `
		INSERT INTO children (id, name, age, neighborhood, reviewed, reviewed_by, reviewed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		ON CONFLICT (id) DO NOTHING`,
		id, name, age, neighborhood, reviewed, reviewedBy,
	)
	if err != nil {
		t.Fatalf("failed to insert child: %v", err)
	}
}

func InsertHealth(t *testing.T, pool *pgxpool.Pool, childID string, vaccinationsUpToDate bool) int64 {
	t.Helper()
	if _, err := pool.Exec(context.Background(), `
		INSERT INTO health (child_id, vaccinations_up_to_date, created_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (child_id) DO NOTHING`,
		childID, vaccinationsUpToDate,
	); err != nil {
		t.Fatalf("failed to insert health: %v", err)
	}
	var id int64
	if err := pool.QueryRow(context.Background(), `SELECT id FROM health WHERE child_id = $1`, childID).Scan(&id); err != nil {
		t.Fatalf("failed to get health id: %v", err)
	}
	return id
}

func InsertHealthAlert(t *testing.T, pool *pgxpool.Pool, healthID int64, code, message string) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `
		INSERT INTO alert_health (health_id, code, message, created_at)
		VALUES ($1, $2, $3, NOW())`,
		healthID, code, message,
	)
	if err != nil {
		t.Fatalf("failed to insert health alert: %v", err)
	}
}

func InsertEducation(t *testing.T, pool *pgxpool.Pool, childID string, schoolName *string, freqPercent int) int64 {
	t.Helper()
	if _, err := pool.Exec(context.Background(), `
		INSERT INTO education (child_id, school_name, frequency_percent, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (child_id) DO NOTHING`,
		childID, schoolName, freqPercent,
	); err != nil {
		t.Fatalf("failed to insert education: %v", err)
	}
	var id int64
	if err := pool.QueryRow(context.Background(), `SELECT id FROM education WHERE child_id = $1`, childID).Scan(&id); err != nil {
		t.Fatalf("failed to get education id: %v", err)
	}
	return id
}

func InsertEducationAlert(t *testing.T, pool *pgxpool.Pool, educationID int64, code, message string) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `
		INSERT INTO alert_education (education_id, code, message, created_at)
		VALUES ($1, $2, $3, NOW())`,
		educationID, code, message,
	)
	if err != nil {
		t.Fatalf("failed to insert education alert: %v", err)
	}
}

func InsertSocialAssistance(t *testing.T, pool *pgxpool.Pool, childID string, cadUnico, activeBenefit bool) int64 {
	t.Helper()
	if _, err := pool.Exec(context.Background(), `
		INSERT INTO social_assistance (child_id, cad_unico, active_benefit, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (child_id) DO NOTHING`,
		childID, cadUnico, activeBenefit,
	); err != nil {
		t.Fatalf("failed to insert social_assistance: %v", err)
	}
	var id int64
	if err := pool.QueryRow(context.Background(), `SELECT id FROM social_assistance WHERE child_id = $1`, childID).Scan(&id); err != nil {
		t.Fatalf("failed to get social_assistance id: %v", err)
	}
	return id
}

func InsertSocialAssistanceAlert(t *testing.T, pool *pgxpool.Pool, socialAssistanceID int64, code, message string) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `
		INSERT INTO alert_social_assistance (social_assistance_id, code, message, created_at)
		VALUES ($1, $2, $3, NOW())`,
		socialAssistanceID, code, message,
	)
	if err != nil {
		t.Fatalf("failed to insert social assistance alert: %v", err)
	}
}
