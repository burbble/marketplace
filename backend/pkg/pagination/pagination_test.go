package pagination

import (
	"testing"
)

func TestPagePaginationGetLimit(t *testing.T) {
	tests := []struct {
		name     string
		pag      PagePagination
		expected uint64
	}{
		{"default", PagePagination{Page: 1}, defaultPageSize},
		{"custom page size", PagePagination{Page: 1, PageSize: 10}, 10},
		{"limit overrides page size", PagePagination{Page: 1, PageSize: 10, Limit: 50}, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pag.GetLimit()
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestPagePaginationGetOffset(t *testing.T) {
	tests := []struct {
		name     string
		pag      PagePagination
		expected uint64
	}{
		{"page 1", PagePagination{Page: 1, PageSize: 24}, 0},
		{"page 2", PagePagination{Page: 2, PageSize: 24}, 24},
		{"page 3 size 10", PagePagination{Page: 3, PageSize: 10}, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pag.GetOffset()
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestSortRequestGetSortClause(t *testing.T) {
	tests := []struct {
		name     string
		req      SortRequest
		expected string
	}{
		{"default", SortRequest{}, "created_at DESC"},
		{"custom field asc", SortRequest{SortBy: "price", SortOrder: "asc"}, "price ASC"},
		{"custom field desc", SortRequest{SortBy: "name", SortOrder: "desc"}, "name DESC"},
		{"invalid order defaults to DESC", SortRequest{SortBy: "price", SortOrder: "invalid"}, "price DESC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.req.GetSortClause()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSortFieldsRequestParseSortFields(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []string
		expectError bool
	}{
		{"empty returns default", "", []string{"created_at DESC"}, false},
		{"single field", "price:asc", []string{"price ASC"}, false},
		{"multiple fields", "price:asc,name:desc", []string{"price ASC", "name DESC"}, false},
		{"default order", "price", []string{"price DESC"}, false},
		{"invalid order", "price:invalid", nil, true},
		{"empty field", ":asc", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sfr := SortFieldsRequest{SortFields: tt.input}
			result, err := sfr.ParseSortFields()

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d clauses, got %d", len(tt.expected), len(result))
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("clause[%d]: expected %q, got %q", i, tt.expected[i], result[i])
				}
			}
		})
	}
}

func TestBuildSortClause(t *testing.T) {
	requests := []SortRequest{
		{SortBy: "price", SortOrder: "asc"},
		{SortBy: "name", SortOrder: "desc"},
	}

	result := BuildSortClause(requests)
	expected := "price ASC, name DESC"

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
