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
		{
			name: "Filter with OR operator - matches either",
			data: [][]string{
				{"Name", "Status", "Age"},
				{"Alice", "Active", "30"},
				{"Bob", "Pending", "25"},
				{"Charlie", "Inactive", "35"},
				{"David", "Active", "40"},
			},
			colIndex:     1, // Status column
			query:        "Active OR Pending",
			expectedRows: 5, // header + Alice (Active), Bob (Pending), Charlie (Inactive has "active"), David (Active)
		},
		{
			name: "Filter with OR operator - case insensitive",
			data: [][]string{
				{"Name", "Status", "Age"},
				{"Alice", "ACTIVE", "30"},
				{"Bob", "pending", "25"},
				{"Charlie", "Inactive", "35"},
			},
			colIndex:     1, // Status column
			query:        "active OR PENDING",
			expectedRows: 4, // header + Alice (ACTIVE), Bob (pending), Charlie (Inactive has "active")
		},
		{
			name: "Filter with AND operator - must have both",
			data: [][]string{
				{"Name", "Description", "Age"},
				{"Alice", "user admin", "30"},
				{"Bob", "user only", "25"},
				{"Charlie", "admin only", "35"},
				{"David", "admin user", "40"},
			},
			colIndex:     1, // Description column
			query:        "user AND admin",
			expectedRows: 3, // header + Alice (user admin), David (admin user)
		},
		{
			name: "Filter with AND operator - no matches",
			data: [][]string{
				{"Name", "Status", "Age"},
				{"Alice", "Active", "30"},
				{"Bob", "Pending", "25"},
				{"Charlie", "Inactive", "35"},
			},
			colIndex:     1, // Status column
			query:        "Active AND Pending",
			expectedRows: 1, // only header (no row has both)
		},
		{
			name: "Filter with OR operator - multiple terms",
			data: [][]string{
				{"Name", "Status", "Age"},
				{"Alice", "red", "30"},
				{"Bob", "blue", "25"},
				{"Charlie", "green", "35"},
				{"David", "yellow", "40"},
			},
			colIndex:     1, // Status column
			query:        "red OR blue OR green",
			expectedRows: 4, // header + Alice, Bob, Charlie
		},
		{
			name: "Filter with AND operator - partial matches",
			data: [][]string{
				{"Name", "Tags", "Age"},
				{"Alice", "python developer", "30"},
				{"Bob", "python", "25"},
				{"Charlie", "developer", "35"},
				{"David", "java developer", "40"},
			},
			colIndex:     1, // Tags column
			query:        "python AND developer",
			expectedRows: 2, // header + Alice (has both)
		},
		{
			name: "Filter with ROR operator - row-level OR",
			data: [][]string{
				{"Name", "Status", "Age"},
				{"Alice", "Active", "30"},
				{"Bob", "Pending", "25"},
				{"Charlie", "Inactive", "35"},
				{"David", "Completed", "40"},
			},
			colIndex:     1, // Status column
			query:        "Active ROR Pending",
			expectedRows: 4, // header + Alice (Active), Bob (Pending), Charlie (Inactive has "active")
		},
		{
			name: "Filter with ROR operator - multiple terms",
			data: [][]string{
				{"Name", "Priority", "Age"},
				{"Alice", "high", "30"},
				{"Bob", "low", "25"},
				{"Charlie", "medium", "35"},
				{"David", "high", "40"},
			},
			colIndex:     1, // Priority column
			query:        "high ROR low",
			expectedRows: 4, // header + Alice (high), Bob (low), David (high)
		},
		{
			name: "Filter with ROR operator - case insensitive",
			data: [][]string{
				{"Name", "Status", "Age"},
				{"Alice", "ACTIVE", "30"},
				{"Bob", "pending", "25"},
				{"Charlie", "Inactive", "35"},
			},
			colIndex:     1, // Status column
			query:        "active ROR PENDING",
			expectedRows: 4, // header + Alice (ACTIVE), Bob (pending), Charlie (Inactive has "active")
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

func TestBuffer_filterByColumn_NumericComparison(t *testing.T) {
	tests := []struct {
		name         string
		data         [][]string
		colIndex     int
		query        string
		expectedRows int // including header
	}{
		{
			name: "Filter > operator",
			data: [][]string{
				{"Name", "Age", "Score"},
				{"Alice", "30", "85"},
				{"Bob", "25", "90"},
				{"Charlie", "35", "78"},
				{"David", "28", "92"},
			},
			colIndex:     1, // Age column
			query:        ">28",
			expectedRows: 3, // header + Alice (30), Charlie (35)
		},
		{
			name: "Filter < operator",
			data: [][]string{
				{"Name", "Age", "Score"},
				{"Alice", "30", "85"},
				{"Bob", "25", "90"},
				{"Charlie", "35", "78"},
				{"David", "28", "92"},
			},
			colIndex:     1, // Age column
			query:        "<30",
			expectedRows: 3, // header + Bob (25), David (28)
		},
		{
			name: "Filter >= operator",
			data: [][]string{
				{"Name", "Age", "Score"},
				{"Alice", "30", "85"},
				{"Bob", "25", "90"},
				{"Charlie", "35", "78"},
				{"David", "28", "92"},
			},
			colIndex:     1, // Age column
			query:        ">=30",
			expectedRows: 3, // header + Alice (30), Charlie (35)
		},
		{
			name: "Filter <= operator",
			data: [][]string{
				{"Name", "Age", "Score"},
				{"Alice", "30", "85"},
				{"Bob", "25", "90"},
				{"Charlie", "35", "78"},
				{"David", "28", "92"},
			},
			colIndex:     1, // Age column
			query:        "<=28",
			expectedRows: 3, // header + Bob (25), David (28)
		},
		{
			name: "Filter > with decimals",
			data: [][]string{
				{"Name", "Score"},
				{"Alice", "85.5"},
				{"Bob", "90.2"},
				{"Charlie", "78.8"},
				{"David", "92.1"},
			},
			colIndex:     1, // Score column
			query:        ">85",
			expectedRows: 4, // header + Alice (85.5), Bob (90.2), David (92.1)
		},
		{
			name: "Filter < with decimals",
			data: [][]string{
				{"Name", "Score"},
				{"Alice", "85.5"},
				{"Bob", "90.2"},
				{"Charlie", "78.8"},
				{"David", "92.1"},
			},
			colIndex:     1, // Score column
			query:        "<85",
			expectedRows: 2, // header + Charlie (78.8)
		},
		{
			name: "Filter > on string column (should not match)",
			data: [][]string{
				{"Name", "Status"},
				{"Alice", "Active"},
				{"Bob", "Inactive"},
				{"Charlie", "Pending"},
			},
			colIndex:     1, // Status column (string type)
			query:        ">5",
			expectedRows: 1, // header only (string column, no numeric comparison)
		},
		{
			name: "Filter > with spaces",
			data: [][]string{
				{"Name", "Age"},
				{"Alice", "30"},
				{"Bob", "25"},
				{"Charlie", "35"},
			},
			colIndex:     1, // Age column
			query:        "> 28",
			expectedRows: 3, // header + Alice (30), Charlie (35)
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
			
			// Detect column types
			b.detectAllColumnTypes()

			// Apply filter
			filtered := b.filterByColumn(tt.colIndex, tt.query, false)

			// Check result
			if filtered.rowLen != tt.expectedRows {
				t.Errorf("Expected %d rows (including header), got %d", tt.expectedRows, filtered.rowLen)
				// Print actual data for debugging
				t.Logf("Actual rows:")
				for i := 0; i < filtered.rowLen; i++ {
					t.Logf("  Row %d: %v", i, filtered.cont[i])
				}
			}

			// Verify header is preserved
			if filtered.rowLen > 0 {
				for col := 0; col < filtered.colLen && col < len(tt.data[0]); col++ {
					if filtered.cont[0][col] != tt.data[0][col] {
						t.Errorf("Header not preserved: expected %s, got %s", tt.data[0][col], filtered.cont[0][col])
					}
				}
			}
		})
	}
}
