package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestMemoryLimit_Basic(t *testing.T) {
	b := createNewBuffer()

	// Check default limit (unlimited = 0)
	if b.maxMemory != 0 {
		t.Errorf("Default memory limit should be 0 (unlimited), got %d", b.maxMemory)
	}

	// Test setting custom limit
	customLimit := int64(100 * 1024 * 1024) // 100 MB
	b.setMemoryLimit(customLimit)

	if b.getMemoryLimit() != customLimit {
		t.Errorf("Memory limit should be %d, got %d", customLimit, b.getMemoryLimit())
	}

	// Test back to unlimited (0)
	b.setMemoryLimit(0)
	if b.getMemoryLimit() != 0 {
		t.Error("Memory limit should be 0 (unlimited)")
	}
}

func TestMemoryLimit_Enforcement(t *testing.T) {
	b := createNewBuffer()

	// Set very low limit to trigger easily
	lowLimit := int64(1024) // 1 KB
	b.setMemoryLimit(lowLimit)

	// Try to add rows that exceed limit
	largeRow := make([]string, 100)
	for i := range largeRow {
		largeRow[i] = strings.Repeat("x", 100) // 100 bytes each
	}

	err := b.contAppendSli(largeRow, false)
	if err == nil {
		t.Error("Expected memory limit error, got nil")
	}

	if !strings.Contains(err.Error(), "Memory limit exceeded") {
		t.Errorf("Expected memory limit error, got: %v", err)
	}
}

func TestMemoryLimit_Tracking(t *testing.T) {
	b := createNewBuffer()

	// Start with 0 usage
	if b.getMemoryUsage() != 0 {
		t.Errorf("Initial memory usage should be 0, got %d", b.getMemoryUsage())
	}

	// Add a row and check usage increases
	testRow := []string{"test1", "test2", "test3"}
	err := b.contAppendSli(testRow, false)
	if err != nil {
		t.Fatalf("Failed to add row: %v", err)
	}

	if b.getMemoryUsage() <= 0 {
		t.Error("Memory usage should increase after adding row")
	}

	initialUsage := b.getMemoryUsage()

	// Add another row
	err = b.contAppendSli(testRow, false)
	if err != nil {
		t.Fatalf("Failed to add second row: %v", err)
	}

	if b.getMemoryUsage() <= initialUsage {
		t.Error("Memory usage should increase after adding second row")
	}
}

func TestMemoryLimit_Stats(t *testing.T) {
	b := createNewBuffer()
	b.setMemoryLimit(1024 * 1024) // 1 MB

	// Add some data
	for i := 0; i < 10; i++ {
		row := []string{
			fmt.Sprintf("value_%d", i),
			fmt.Sprintf("data_%d", i),
			fmt.Sprintf("test_%d", i),
		}
		err := b.contAppendSli(row, false)
		if err != nil {
			t.Fatalf("Failed to add row: %v", err)
		}
	}

	stats := b.getMemoryStats()

	// Check required fields
	if _, ok := stats["current_bytes"]; !ok {
		t.Error("Stats should include current_bytes")
	}

	if _, ok := stats["current_formatted"]; !ok {
		t.Error("Stats should include current_formatted")
	}

	if _, ok := stats["limit_bytes"]; !ok {
		t.Error("Stats should include limit_bytes")
	}

	if _, ok := stats["limit_formatted"]; !ok {
		t.Error("Stats should include limit_formatted")
	}

	if _, ok := stats["usage_percent"]; !ok {
		t.Error("Stats should include usage_percent")
	}

	// Check values
	currentBytes := stats["current_bytes"].(int64)
	if currentBytes <= 0 {
		t.Error("Current bytes should be > 0")
	}

	limitBytes := stats["limit_bytes"].(int64)
	if limitBytes != 1024*1024 {
		t.Errorf("Limit should be 1048576, got %d", limitBytes)
	}

	usagePercent := stats["usage_percent"].(float64)
	if usagePercent < 0 || usagePercent > 100 {
		t.Errorf("Usage percent should be 0-100, got %f", usagePercent)
	}

	t.Logf("Memory stats: %+v", stats)
}

func TestMemoryLimit_UnlimitedMode(t *testing.T) {
	b := createNewBuffer()
	b.setMemoryLimit(0) // Unlimited

	// Should be able to add many rows without limit
	for i := 0; i < 1000; i++ {
		row := []string{
			strings.Repeat("x", 100),
			strings.Repeat("y", 100),
			strings.Repeat("z", 100),
		}
		err := b.contAppendSli(row, false)
		if err != nil {
			t.Fatalf("Should not fail in unlimited mode: %v", err)
		}
	}

	stats := b.getMemoryStats()
	if unlimited, ok := stats["unlimited"]; !ok || !unlimited.(bool) {
		t.Error("Stats should indicate unlimited mode")
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{100, "100 B"},
		{1023, "1023 B"},
		{1024, "1.00 KB"},
		{1024 * 1024, "1.00 MB"},
		{1024 * 1024 * 1024, "1.00 GB"},
		{1536 * 1024 * 1024, "1.50 GB"},
		{2 * 1024 * 1024 * 1024, "2.00 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatBytes(tt.bytes)
			if result != tt.expected {
				t.Errorf("formatBytes(%d) = %s, want %s", tt.bytes, result, tt.expected)
			}
		})
	}
}

func TestMemoryLimit_EstimateRowSize(t *testing.T) {
	b := createNewBuffer()

	testCases := []struct {
		name string
		row  []string
	}{
		{"Empty row", []string{}},
		{"Single cell", []string{"test"}},
		{"Multiple cells", []string{"cell1", "cell2", "cell3"}},
		{"Large cells", []string{strings.Repeat("x", 1000), strings.Repeat("y", 1000)}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			size := b.estimateRowSize(tc.row)
			if size < 0 {
				t.Error("Row size should not be negative")
			}

			// Rough validation - size should be reasonable
			minSize := int64(len(tc.row) * 8) // At least pointer overhead
			if size < minSize {
				t.Errorf("Estimated size %d is too small (min %d)", size, minSize)
			}

			t.Logf("Row with %d cells: estimated %d bytes", len(tc.row), size)
		})
	}
}

func TestMemoryLimit_GracefulDegradation(t *testing.T) {
	b := createNewBuffer()

	// Set moderate limit
	limit := int64(10 * 1024) // 10 KB
	b.setMemoryLimit(limit)

	// Add rows until we hit limit
	rowsAdded := 0
	for i := 0; i < 1000; i++ {
		row := []string{
			fmt.Sprintf("data_%d", i),
			strings.Repeat("x", 50),
			strings.Repeat("y", 50),
		}

		err := b.contAppendSli(row, false)
		if err != nil {
			// Should fail gracefully with memory limit error
			if !strings.Contains(err.Error(), "Memory limit exceeded") {
				t.Errorf("Expected memory limit error, got: %v", err)
			}
			break
		}
		rowsAdded++
	}

	if rowsAdded == 0 {
		t.Error("Should have been able to add at least one row")
	}

	if rowsAdded >= 1000 {
		t.Error("Should have hit memory limit before 1000 rows")
	}

	t.Logf("Successfully added %d rows before hitting %s limit", rowsAdded, formatBytes(limit))

	// Verify buffer is still usable
	if b.rowLen != rowsAdded {
		t.Errorf("Row count mismatch: %d != %d", b.rowLen, rowsAdded)
	}
}

func TestMemoryLimit_Integration(t *testing.T) {
	// Simulate real-world scenario
	data := [][]string{
		{"Name", "Email", "Country", "Status"},
	}

	// Add 1000 rows of realistic data
	for i := 0; i < 1000; i++ {
		data = append(data, []string{
			fmt.Sprintf("User_%d", i),
			fmt.Sprintf("user%d@example.com", i),
			"USA",
			"Active",
		})
	}

	// Test with reasonable limit
	b, err := createNewBufferWithData(data, false)
	if err != nil {
		t.Fatalf("Failed to create buffer: %v", err)
	}

	usage := b.getMemoryUsage()
	if usage <= 0 {
		t.Error("Memory usage should be tracked")
	}

	t.Logf("Buffer with 1000 rows uses approximately %s", formatBytes(usage))

	// Verify it's within reasonable bounds (should be < 1MB for this data)
	if usage > 1024*1024 {
		t.Errorf("Memory usage seems too high: %s", formatBytes(usage))
	}
}

// Benchmark memory tracking overhead
func BenchmarkMemoryTracking(b *testing.B) {
	buf := createNewBuffer()
	buf.setMemoryLimit(0) // Unlimited to avoid limit checks

	row := []string{"test1", "test2", "test3", "test4", "test5"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = buf.estimateRowSize(row)
	}
}

func BenchmarkMemoryLimitCheck(b *testing.B) {
	buf := createNewBuffer()
	buf.setMemoryLimit(1024 * 1024 * 1024) // 1GB

	row := []string{"test1", "test2", "test3"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = buf.contAppendSli(row, false)
	}
}
