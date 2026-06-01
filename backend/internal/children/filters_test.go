package children

import (
	"testing"
)

func TestFilters_Normalize_DefaultPage(t *testing.T) {
	f := Filters{Page: 0, PerPage: 10}
	f.Normalize()
	if f.Page != 1 {
		t.Errorf("expected page 1, got %d", f.Page)
	}
}

func TestFilters_Normalize_MinPerPage(t *testing.T) {
	f := Filters{Page: 1, PerPage: 5}
	f.Normalize()
	if f.PerPage != 10 {
		t.Errorf("expected per_page 10, got %d", f.PerPage)
	}
}

func TestFilters_Normalize_MaxPerPage(t *testing.T) {
	f := Filters{Page: 1, PerPage: 100}
	f.Normalize()
	if f.PerPage != 50 {
		t.Errorf("expected per_page 50, got %d", f.PerPage)
	}
}

func TestFilters_Normalize_ValidValues(t *testing.T) {
	f := Filters{Page: 3, PerPage: 25}
	f.Normalize()
	if f.Page != 3 {
		t.Errorf("expected page 3, got %d", f.Page)
	}
	if f.PerPage != 25 {
		t.Errorf("expected per_page 25, got %d", f.PerPage)
	}
}

func TestFilters_Normalize_NegativePage(t *testing.T) {
	f := Filters{Page: -5, PerPage: 10}
	f.Normalize()
	if f.Page != 1 {
		t.Errorf("expected page 1, got %d", f.Page)
	}
}

func TestFilters_Normalize_ZeroValues(t *testing.T) {
	f := Filters{}
	f.Normalize()
	if f.Page != 1 {
		t.Errorf("expected page 1, got %d", f.Page)
	}
	if f.PerPage != 10 {
		t.Errorf("expected per_page 10, got %d", f.PerPage)
	}
}

func TestFilters_Offset_FirstPage(t *testing.T) {
	f := Filters{Page: 1, PerPage: 10}
	f.Normalize()
	if off := f.Offset(); off != 0 {
		t.Errorf("expected offset 0, got %d", off)
	}
}

func TestFilters_Offset_ThirdPage(t *testing.T) {
	f := Filters{Page: 3, PerPage: 10}
	f.Normalize()
	if off := f.Offset(); off != 20 {
		t.Errorf("expected offset 20, got %d", off)
	}
}

func TestFilters_Offset_CustomPerPage(t *testing.T) {
	f := Filters{Page: 2, PerPage: 25}
	f.Normalize()
	if off := f.Offset(); off != 25 {
		t.Errorf("expected offset 25, got %d", off)
	}
}

func TestFilters_Normalize_MinPerPageBoundary(t *testing.T) {
	f := Filters{Page: 1, PerPage: 10}
	f.Normalize()
	if f.PerPage != 10 {
		t.Errorf("expected per_page 10, got %d", f.PerPage)
	}
}

func TestFilters_Normalize_MaxPerPageBoundary(t *testing.T) {
	f := Filters{Page: 1, PerPage: 50}
	f.Normalize()
	if f.PerPage != 50 {
		t.Errorf("expected per_page 50, got %d", f.PerPage)
	}
}
