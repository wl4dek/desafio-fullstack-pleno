package database_test

import (
	"context"
	"testing"

	"backend/internal/database/testutil"
)

func TestIntegration_Connect(t *testing.T) {
	testutil.SkipIfNoDB(t)

	pool := testutil.Connect(t)
	if pool == nil {
		t.Fatal("expected non-nil pool")
	}

	err := pool.Ping(context.Background())
	if err != nil {
		t.Fatalf("failed to ping: %v", err)
	}
}

func TestIntegration_RunMigrations(t *testing.T) {
	testutil.SkipIfNoDB(t)

	dsn := testutil.DSN(t)
	pool := testutil.Connect(t)
	defer pool.Close()

	migrationsPath := testutil.MigrationsPath(t)

	testutil.Migrate(t, dsn, migrationsPath)

	var tableCount int
	err := pool.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM information_schema.tables
		WHERE table_schema = 'public'`).Scan(&tableCount)
	if err != nil {
		t.Fatalf("failed to count tables: %v", err)
	}
	if tableCount < 6 {
		t.Errorf("expected at least 6 tables, got %d", tableCount)
	}
}

func TestIntegration_LoadSeed(t *testing.T) {
	testutil.SkipIfNoDB(t)

	dsn := testutil.DSN(t)
	pool := testutil.Connect(t)
	defer pool.Close()

	migrationsPath := testutil.MigrationsPath(t)
	seedPath := testutil.SeedPath(t)

	testutil.Migrate(t, dsn, migrationsPath)
	testutil.Seed(t, pool, seedPath)

	var count int
	err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM children").Scan(&count)
	if err != nil {
		t.Fatalf("failed to count children: %v", err)
	}
	if count == 0 {
		t.Fatal("expected seed data to be loaded")
	}
}

func TestIntegration_LoadSeed_Idempotent(t *testing.T) {
	testutil.SkipIfNoDB(t)

	dsn := testutil.DSN(t)
	pool := testutil.Connect(t)
	defer pool.Close()

	migrationsPath := testutil.MigrationsPath(t)
	seedPath := testutil.SeedPath(t)

	testutil.Migrate(t, dsn, migrationsPath)

	testutil.Seed(t, pool, seedPath)

	var count1 int
	pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM children").Scan(&count1)

	testutil.Seed(t, pool, seedPath)

	var count2 int
	pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM children").Scan(&count2)

	if count1 != count2 {
		t.Errorf("seed should be idempotent: before %d, after %d", count1, count2)
	}
}

func TestIntegration_RunMigrations_Idempotent(t *testing.T) {
	testutil.SkipIfNoDB(t)

	dsn := testutil.DSN(t)
	pool := testutil.Connect(t)
	defer pool.Close()

	migrationsPath := testutil.MigrationsPath(t)

	testutil.Migrate(t, dsn, migrationsPath)
	testutil.Migrate(t, dsn, migrationsPath)

	var tableCount int
	err := pool.QueryRow(context.Background(), `
		SELECT COUNT(*) FROM information_schema.tables
		WHERE table_schema = 'public'`).Scan(&tableCount)
	if err != nil {
		t.Fatalf("failed to count tables: %v", err)
	}
	if tableCount < 6 {
		t.Errorf("expected at least 6 tables, got %d", tableCount)
	}
}
