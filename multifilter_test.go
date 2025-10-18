package main

import (
	"testing"
)

func TestMultiColumnFilter(t *testing.T) {
	// Create test data
	data := [][]string{
		{"Name", "Age", "City", "Department"},
		{"Alice", "30", "New York", "Engineering"},
		{"Bob", "25", "San Francisco", "Sales"},
		{"Charlie", "35", "New York", "Engineering"},
		{"David", "28", "Chicago", "Marketing"},
		{"Eve", "32", "New York", "Sales"},
		{"Frank", "45", "San Francisco", "Engineering"},
	}

	b, err := createNewBufferWithData(data, true)
	if err != nil {
		t.Fatalf("Failed to create buffer: %v", err)
	}
	b.rowFreeze = 1

	t.Run("Apply single filter", func(t *testing.T) {
		// Filter by City = "New York"
		filtered := b.filterByColumn(2, "New York", false)
		
		// Should have header + 3 data rows
		expected := 4 // 1 header + 3 matching rows
		if filtered.rowLen != expected {
			t.Errorf("Expected %d rows, got %d", expected, filtered.rowLen)
		}
	})

	t.Run("Apply multiple filters sequentially", func(t *testing.T) {
		// First filter by City = "New York"
		filtered1 := b.filterByColumn(2, "New York", false)
		
		// Then filter by Department = "Engineering"
		filtered2 := filtered1.filterByColumn(3, "Engineering", false)
		
		// Should have header + 2 data rows (Alice and Charlie)
		expected := 3 // 1 header + 2 matching rows
		if filtered2.rowLen != expected {
			t.Errorf("Expected %d rows, got %d", expected, filtered2.rowLen)
		}
		
		// Verify the correct rows remain
		if filtered2.cont[1][0] != "Alice" {
			t.Errorf("Expected Alice in row 1, got %s", filtered2.cont[1][0])
		}
		if filtered2.cont[2][0] != "Charlie" {
			t.Errorf("Expected Charlie in row 2, got %s", filtered2.cont[2][0])
		}
	})

	t.Run("Apply three filters", func(t *testing.T) {
		// Filter by City = "New York"
		filtered1 := b.filterByColumn(2, "New York", false)
		
		// Filter by Department containing "Eng"
		filtered2 := filtered1.filterByColumn(3, "Eng", false)
		
		// Filter by Age > 32 (filter for "3" to match 30 and 35)
		filtered3 := filtered2.filterByColumn(1, "3", false)
		
		// Should have header + 2 data rows (Alice:30, Charlie:35)
		expected := 3
		if filtered3.rowLen != expected {
			t.Errorf("Expected %d rows after 3 filters, got %d", expected, filtered3.rowLen)
		}
	})

	t.Run("Filter with no results", func(t *testing.T) {
		// Filter by City = "New York"
		filtered1 := b.filterByColumn(2, "New York", false)
		
		// Filter by Department = "Marketing" (no New York Marketing employees)
		filtered2 := filtered1.filterByColumn(3, "Marketing", false)
		
		// Should have header only
		expected := 1
		if filtered2.rowLen != expected {
			t.Errorf("Expected %d row (header only), got %d", expected, filtered2.rowLen)
		}
	})

	t.Run("Filters preserve column types", func(t *testing.T) {
		b.colType[1] = colTypeFloat // Set Age as numeric
		b.colType[2] = colTypeStr   // Set City as string
		
		filtered := b.filterByColumn(2, "New York", false)
		
		// Check that column types are preserved
		if filtered.colType[1] != colTypeFloat {
			t.Errorf("Column type for Age not preserved after filter")
		}
		if filtered.colType[2] != colTypeStr {
			t.Errorf("Column type for City not preserved after filter")
		}
	})

	t.Run("Multiple filters with operators", func(t *testing.T) {
		// Filter by City with OR operator
		filtered1 := b.filterByColumn(2, "New York OR Chicago", false)
		
		// Then filter by Department
		filtered2 := filtered1.filterByColumn(3, "Engineering OR Marketing", false)
		
		// Should match: Alice, Charlie (NY Eng), David (Chicago Marketing)
		expected := 4 // 1 header + 3 matching rows
		if filtered2.rowLen != expected {
			t.Errorf("Expected %d rows with OR filters, got %d", expected, filtered2.rowLen)
		}
	})

	t.Run("Remove one filter from multiple active filters", func(t *testing.T) {
		// Apply two filters
		filtered1 := b.filterByColumn(2, "New York", false) // City filter
		filtered2 := filtered1.filterByColumn(3, "Engineering", false) // Department filter
		
		// Should have 3 rows (header + Alice + Charlie)
		if filtered2.rowLen != 3 {
			t.Errorf("Expected 3 rows with both filters, got %d", filtered2.rowLen)
		}
		
		// Now remove the Department filter by applying City filter to original
		// This simulates what happens when user removes one filter
		filteredAfterRemoval := b.filterByColumn(2, "New York", false)
		
		// Should have 4 rows (header + Alice + Charlie + Eve)
		expected := 4
		if filteredAfterRemoval.rowLen != expected {
			t.Errorf("Expected %d rows after removing one filter, got %d", expected, filteredAfterRemoval.rowLen)
		}
	})
}
