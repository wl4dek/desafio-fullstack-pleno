package children

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewQueryBuilder(t *testing.T) {
	qb := NewQueryBuilder()
	if qb == nil {
		t.Fatal("expected non-nil QueryBuilder")
	}
	if qb.baseQuery == "" {
		t.Error("expected non-empty baseQuery")
	}
	if qb.countQuery == "" {
		t.Error("expected non-empty countQuery")
	}
	if qb.findByIDQuery == "" {
		t.Error("expected non-empty findByIDQuery")
	}
}

func TestQueryBuilder_AddCondition(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddCondition("c.name ILIKE $1", "%test%")

	if len(qb.where) != 1 {
		t.Errorf("expected 1 condition, got %d", len(qb.where))
	}
	if len(qb.args) != 1 {
		t.Errorf("expected 1 arg, got %d", len(qb.args))
	}
	if qb.args[0] != "%test%" {
		t.Errorf("expected arg %%test%%, got %v", qb.args[0])
	}
}

func TestQueryBuilder_AddCondition_Multiple(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddCondition("c.name ILIKE $1", "%test%")
	qb.AddCondition("c.neighborhood = $2", "Centro")

	if len(qb.where) != 2 {
		t.Errorf("expected 2 conditions, got %d", len(qb.where))
	}
	if len(qb.args) != 2 {
		t.Errorf("expected 2 args, got %d", len(qb.args))
	}
}

func TestQueryBuilder_AddConditionOnly(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddConditionOnly("c.reviewed = true")

	if len(qb.where) != 1 {
		t.Errorf("expected 1 condition, got %d", len(qb.where))
	}
	if len(qb.args) != 0 {
		t.Errorf("expected 0 args, got %d", len(qb.args))
	}
}

func TestQueryBuilder_BuildList_NoFilters(t *testing.T) {
	qb := NewQueryBuilder()
	query, args := qb.BuildList()

	if !strings.Contains(query, "FROM children c") {
		t.Error("expected FROM children c in query")
	}
	if !strings.Contains(query, "ORDER BY c.created_at DESC") {
		t.Error("expected ORDER BY c.created_at DESC in query")
	}
	countWhere := strings.Count(query, "WHERE") - strings.Count(query, "THEN")
	if countWhere != 0 {
		t.Errorf("expected no top-level WHERE clause, but query has %d WHERE in SELECT", countWhere)
	}
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestQueryBuilder_BuildList_WithFilters(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddCondition("c.neighborhood = $1", "Centro")
	query, args := qb.BuildList()

	if !strings.Contains(query, "WHERE") {
		t.Error("expected WHERE clause")
	}
	if !strings.Contains(query, "c.neighborhood = $1") {
		t.Error("expected neighborhood filter in query")
	}
	if len(args) != 1 {
		t.Errorf("expected 1 arg, got %d", len(args))
	}
}

func TestQueryBuilder_BuildPaginatedList(t *testing.T) {
	qb := NewQueryBuilder()
	query, _ := qb.BuildPaginatedList(10, 0)

	if !strings.Contains(query, "LIMIT 10") {
		t.Error("expected LIMIT 10 in query")
	}
	if !strings.Contains(query, "OFFSET 0") {
		t.Error("expected OFFSET 0 in query")
	}
}

func TestQueryBuilder_BuildPaginatedList_WithOffset(t *testing.T) {
	qb := NewQueryBuilder()
	query, _ := qb.BuildPaginatedList(25, 50)

	if !strings.Contains(query, "LIMIT 25") {
		t.Error("expected LIMIT 25 in query")
	}
	if !strings.Contains(query, "OFFSET 50") {
		t.Error("expected OFFSET 50 in query")
	}
}

func TestQueryBuilder_BuildCount_NoFilters(t *testing.T) {
	qb := NewQueryBuilder()
	query, args := qb.BuildCount()

	if !strings.Contains(query, "SELECT COUNT(*) FROM children c") {
		t.Error("expected SELECT COUNT(*) FROM children c")
	}
	if strings.Contains(query, "WHERE") {
		t.Error("expected no WHERE clause without filters")
	}
	if len(args) != 0 {
		t.Errorf("expected 0 args, got %d", len(args))
	}
}

func TestQueryBuilder_BuildCount_WithFilters(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddCondition("c.neighborhood = $1", "Centro")
	query, _ := qb.BuildCount()

	if !strings.Contains(query, "WHERE") {
		t.Error("expected WHERE clause")
	}
	if !strings.Contains(query, "c.neighborhood = $1") {
		t.Error("expected neighborhood filter")
	}
}

func TestQueryBuilder_BuildById(t *testing.T) {
	qb := NewQueryBuilder()
	query, args := qb.BuildById("child-123")

	if !strings.Contains(query, "WHERE c.id = $1") {
		t.Error("expected WHERE c.id = $1 in query")
	}
	if len(args) != 1 {
		t.Errorf("expected 1 arg, got %d", len(args))
	}
	if args[0] != "child-123" {
		t.Errorf("expected child-123, got %v", args[0])
	}
}

func TestQueryBuilder_BuildById_IncludesJoins(t *testing.T) {
	qb := NewQueryBuilder()
	query, _ := qb.BuildById("test")

	expectedJoins := []string{
		"LEFT JOIN health",
		"LEFT JOIN education",
		"LEFT JOIN social_assistance",
	}
	for _, join := range expectedJoins {
		if !strings.Contains(query, join) {
			t.Errorf("expected %s in query", join)
		}
	}
}

func TestQueryBuilder_BuildList_IncludesAlertCategories(t *testing.T) {
	qb := NewQueryBuilder()
	query, _ := qb.BuildList()

	expectedFields := []string{
		"ARRAY_REMOVE",
		"alert_categories",
		"FROM children c",
	}
	for _, f := range expectedFields {
		if !strings.Contains(query, f) {
			t.Errorf("expected %s in query", f)
		}
	}
}

func TestQueryBuilder_BuildPaginatedList_PreservesOrder(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddCondition("c.reviewed = $1", true)
	query, _ := qb.BuildPaginatedList(10, 0)

	orderIdx := strings.Index(query, "ORDER BY")
	limitIdx := strings.Index(query, "LIMIT")

	if orderIdx < 0 {
		t.Fatal("expected ORDER BY in query")
	}
	if limitIdx < 0 {
		t.Fatal("expected LIMIT in query")
	}
	if orderIdx > limitIdx {
		t.Error("expected ORDER BY before LIMIT")
	}
}

func TestQueryBuilder_BuildCount_NoWhereWithEmptyConditions(t *testing.T) {
	qb := NewQueryBuilder()
	query, _ := qb.BuildCount()

	if strings.Contains(query, "WHERE") {
		t.Error("expected no WHERE clause with no conditions")
	}
}

func TestQueryBuilder_MultipleFilters(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddCondition("c.neighborhood = $1", "Centro")
	qb.AddCondition("c.reviewed = $2", true)

	query, _ := qb.BuildList()
	if !strings.Contains(query, "c.neighborhood = $1") {
		t.Error("expected neighborhood filter")
	}
	if !strings.Contains(query, "c.reviewed = $2") {
		t.Error("expected reviewed filter")
	}
}

func TestQueryBuilder_ArgsArePassedThrough(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddCondition("c.name ILIKE $1", "%search%")
	qb.AddCondition("c.neighborhood = $2", "North")

	_, args := qb.BuildList()
	if len(args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(args))
	}
	if args[0] != "%search%" {
		t.Errorf("expected %%search%%, got %v", args[0])
	}
	if args[1] != "North" {
		t.Errorf("expected North, got %v", args[1])
	}
}

func TestQueryBuilder_CountArgsMatchFilters(t *testing.T) {
	qb := NewQueryBuilder()
	qb.AddCondition("c.neighborhood = $1", "Centro")

	_, listArgs := qb.BuildList()
	_, countArgs := qb.BuildCount()

	if len(listArgs) != len(countArgs) {
		t.Errorf("list args (%d) and count args (%d) should match", len(listArgs), len(countArgs))
	}
}

func TestQueryBuilder_BuildById_NoExtraConditions(t *testing.T) {
	qb := NewQueryBuilder()
	query, _ := qb.BuildById("test")

	if strings.Contains(query, "ORDER BY") {
		t.Error("findByID should not contain ORDER BY")
	}
	if strings.Contains(query, "LIMIT") {
		t.Error("findByID should not contain LIMIT")
	}
}

func TestQueryBuilder_AddCondition_DollarSignIndexing(t *testing.T) {
	qb := NewQueryBuilder()

	placeholder := fmt.Sprintf("c.neighborhood = $%d", len(qb.args)+1)
	qb.AddCondition(placeholder, "Centro")

	if qb.where[0] != "c.neighborhood = $1" {
		t.Errorf("expected c.neighborhood = $1, got %s", qb.where[0])
	}
}
