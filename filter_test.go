package main

import (
	"testing"
)

func TestBuffer_filterByColumn(t *testing.T) {
	tests := []struct {
		name          string
		data          [][]string
		colIndex      int
		query         string
		caseSensitive bool
		expectedRows  int // including header
	}{
		{
			name: "Filter by exact match",
			data: [][]string{
				{"Name", "Status", "Age"},
				{"Alice", "Active", "30"},
				{"Bob", "Inactive", "25"},
				{"Charlie", "Active", "35"},
			},
			colIndex:     1, // Status column
			query:        "Active",
			expectedRows: 4, // header + 2 Active rows + 1 Inactive (contains "active")
		},
		{
			name: "Filter by partial match",
			data: [][]string{
				{"Name", "Status", "Age"},
				{"Alice", "Active", "30"},
				{"Bob", "Inactive", "25"},
				{"Charlie", "Active", "35"},
			},
			colIndex:     1, // Status column
			query:        "act",
			expectedRows: 4, // header + all 3 rows containing "act"
		},
		{
			name: "Filter case insensitive",
			data: [][]string{
				{"Name", "Status", "Age"},
				{"Alice", "ACTIVE", "30"},
				{"Bob", "inactive", "25"},
				{"Charlie", "Active", "35"},
			},
			colIndex:     1, // Status column
			query:        "active",
			expectedRows: 4, // header + all 3 rows (case insensitive)
		},
		{
			name: "Filter no matches",
			data: [][]string{
				{"Name", "Status", "Age"},
				{"Alice", "Active", "30"},
				{"Bob", "Inactive", "25"},
				{"Charlie", "Active", "35"},
			},
			colIndex:     1, // Status column
			query:        "Pending",
			expectedRows: 1, // only header
		},
		{
			name: "Filter first column",
			data: [][]string{
				{"Name", "Status", "Age"},
				{"Alice", "Active", "30"},
				{"Bob", "Inactive", "25"},
				{"Charlie", "Active", "35"},
			},
			colIndex:     0, // Name column
			query:        "li",
			expectedRows: 3, // header + Alice, Charlie (both contain "li")
		},
		{
			name: "Empty filter query matches all",
			data: [][]string{
				{"Name", "Status", "Age"},
				{"Alice", "Active", "30"},
				{"Bob", "Inactive", "25"},
			},
			colIndex:     1, // Status column
			query:        "",
			expectedRows: 3, // header + all data rows
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create buffer with test data
			b, err := createNewBufferWithData(tt.data, false)
			if err != nil {
				t.Fatalf("Failed to create buffer: %v", err)
			}
			b.rowFreeze = 1 // Set header row

			// Apply filter
			filtered := b.filterByColumn(tt.colIndex, tt.query, tt.caseSensitive)

			// Check result
			if filtered.rowLen != tt.expectedRows {
				t.Errorf("Expected %d rows (including header), got %d", tt.expectedRows, filtered.rowLen)
			}

			// Verify header is preserved
			if filtered.rowLen > 0 {
				for col := 0; col < filtered.colLen && col < len(tt.data[0]); col++ {
					if filtered.cont[0][col] != tt.data[0][col] {
						t.Errorf("Header not preserved: expected %s, got %s", tt.data[0][col], filtered.cont[0][col])
					}
				}
			}

			// Verify freeze settings are preserved
			if filtered.rowFreeze != b.rowFreeze {
				t.Errorf("rowFreeze not preserved: expected %d, got %d", b.rowFreeze, filtered.rowFreeze)
			}
			if filtered.colFreeze != b.colFreeze {
				t.Errorf("colFreeze not preserved: expected %d, got %d", b.colFreeze, filtered.colFreeze)
			}
		})
	}
}

func TestBuffer_filterByColumn_EdgeCases(t *testing.T) {
	t.Run("Filter column out of bounds", func(t *testing.T) {
		data := [][]string{
			{"Name", "Status"},
			{"Alice", "Active"},
		}
		b, _ := createNewBufferWithData(data, false)
		b.rowFreeze = 1

		// Filter by non-existent column
		filtered := b.filterByColumn(10, "test", false)

		// Should return only header
		if filtered.rowLen != 1 {
			t.Errorf("Expected 1 row (header only), got %d", filtered.rowLen)
		}
	})

	t.Run("Filter empty buffer", func(t *testing.T) {
		b := createNewBuffer()
		b.rowFreeze = 0

		// Filter empty buffer
		filtered := b.filterByColumn(0, "test", false)

		// Should return empty buffer
		if filtered.rowLen != 0 {
			t.Errorf("Expected 0 rows, got %d", filtered.rowLen)
		}
	})

	t.Run("Filter buffer with no header", func(t *testing.T) {
		data := [][]string{
			{"Alice", "Active"},
			{"Bob", "Inactive"},
		}
		b, _ := createNewBufferWithData(data, false)
		b.rowFreeze = 0 // No header

		// Filter by column
		filtered := b.filterByColumn(1, "Active", false)

		// Should return 2 rows (both contain "active")
		if filtered.rowLen != 2 {
			t.Errorf("Expected 2 rows, got %d", filtered.rowLen)
		}
	})
}
