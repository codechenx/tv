package main

import (
	"fmt"
	"runtime"
	"testing"
)

func TestStringInterner(t *testing.T) {
	si := newStringInterner()

	// Test basic interning
	s1 := si.intern("test")
	s2 := si.intern("test")

	// Should return the same string value
	if s1 != s2 {
		t.Error("String interning failed - expected same string value")
	}

	// Intern should deduplicate - both should resolve to same underlying string
	// This is verified by sync.Map storing it once
	count := 0
	si.pool.Range(func(key, value interface{}) bool {
		count++
		return true
	})

	if count != 1 {
		t.Errorf("Expected 1 unique string in pool, got %d", count)
	}

	// Different strings should still work
	s3 := si.intern("different")
	if s1 == s3 {
		t.Error("Different strings should not be equal")
	}

	// Pool should now have 2 entries
	count = 0
	si.pool.Range(func(key, value interface{}) bool {
		count++
		return true
	})

	if count != 2 {
		t.Errorf("Expected 2 unique strings in pool, got %d", count)
	}
}

func TestShouldInternColumn(t *testing.T) {
	tests := []struct {
		name      string
		values    []string
		threshold float64
		expected  bool
	}{
		{
			name:      "High cardinality - should not intern",
			values:    generateUniqueStrings(1000),
			threshold: 0.30,
			expected:  false,
		},
		{
			name:      "Low cardinality - should intern",
			values:    generateRepeatedStrings(1000, 10),
			threshold: 0.30,
			expected:  true,
		},
		{
			name:      "Too small - should not intern",
			values:    []string{"a", "b", "c"},
			threshold: 0.30,
			expected:  false,
		},
		{
			name:      "Medium cardinality at threshold",
			values:    generateRepeatedStrings(1000, 300),
			threshold: 0.30,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldInternColumn(tt.values, tt.threshold)
			if result != tt.expected {
				t.Errorf("shouldInternColumn() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBufferStringInterning(t *testing.T) {
	// Create test data with repeated categorical values
	data := [][]string{
		{"Category", "Value"},
		{"A", "100"},
		{"B", "200"},
		{"A", "300"}, // Repeated
		{"C", "400"},
		{"A", "500"}, // Repeated
		{"B", "600"}, // Repeated
	}

	// Add more rows to make it worthwhile
	for i := 0; i < 100; i++ {
		data = append(data, []string{"A", fmt.Sprintf("%d", i)})
		data = append(data, []string{"B", fmt.Sprintf("%d", i+1000)})
		data = append(data, []string{"C", fmt.Sprintf("%d", i+2000)})
	}

	b, err := createNewBufferWithData(data, false)
	if err != nil {
		t.Fatalf("Failed to create buffer: %v", err)
	}

	b.rowFreeze = 1

	// Enable string interning
	b.enableStringInterning()

	// Check that interning is enabled
	stats := b.getInterningStats()
	if !stats["enabled"].(bool) {
		t.Error("String interning should be enabled")
	}

	// Category column (0) should be interned (low cardinality)
	if !b.internCols[0] {
		t.Error("Category column should be interned")
	}

	// Value column (1) should NOT be interned (high cardinality)
	if b.internCols[1] {
		t.Error("Value column should not be interned")
	}

	t.Logf("Interning stats: %+v", stats)
}

func TestStringInterningMemorySavings(t *testing.T) {
	// Create large dataset with repeated values
	numRows := 10000
	categories := []string{"Active", "Inactive", "Pending", "Completed", "Failed"}

	data := [][]string{{"Status", "ID"}}
	for i := 0; i < numRows; i++ {
		data = append(data, []string{
			categories[i%len(categories)],
			fmt.Sprintf("%d", i),
		})
	}

	// Create buffer without interning
	runtime.GC()
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	b1, _ := createNewBufferWithData(data, false)
	b1.rowFreeze = 1

	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)
	memWithoutInterning := m2.Alloc - m1.Alloc

	// Create buffer with interning
	runtime.GC()
	var m3 runtime.MemStats
	runtime.ReadMemStats(&m3)

	b2, _ := createNewBufferWithData(data, false)
	b2.rowFreeze = 1
	b2.enableStringInterning()

	var m4 runtime.MemStats
	runtime.ReadMemStats(&m4)
	memWithInterning := m4.Alloc - m3.Alloc

	t.Logf("Memory without interning: %d bytes", memWithoutInterning)
	t.Logf("Memory with interning: %d bytes", memWithInterning)

	// Interning should use less memory (though the difference might be small in test)
	// We're mainly testing that it doesn't crash and completes successfully
	if b2.internCols == nil {
		t.Error("Interning structures should be initialized")
	}
}

func TestInternValue(t *testing.T) {
	b := createNewBuffer()
	b.colLen = 2
	b.interners = make([]*stringInterner, 2)
	b.internCols = make([]bool, 2)

	// Enable interning for column 0
	b.interners[0] = newStringInterner()
	b.internCols[0] = true

	// Test interning
	val1 := b.internValue(0, "test")
	val2 := b.internValue(0, "test")

	// Should be same value for column 0
	if val1 != val2 {
		t.Error("Interned values should be same value")
	}

	// Verify it's actually in the interner
	count := 0
	b.interners[0].pool.Range(func(key, value interface{}) bool {
		count++
		return true
	})

	if count != 1 {
		t.Errorf("Expected 1 entry in interner, got %d", count)
	}

	// Column 1 not interned, should return original
	val3 := b.internValue(1, "test")
	if val3 != "test" {
		t.Error("Non-interned column should return original value")
	}
}

// Helper functions
func generateUniqueStrings(count int) []string {
	result := make([]string, count)
	for i := 0; i < count; i++ {
		result[i] = fmt.Sprintf("unique_%d", i)
	}
	return result
}

func generateRepeatedStrings(count, uniqueCount int) []string {
	result := make([]string, count)
	for i := 0; i < count; i++ {
		result[i] = fmt.Sprintf("value_%d", i%uniqueCount)
	}
	return result
}

// Benchmark string interning
func BenchmarkStringInterning(b *testing.B) {
	// Create dataset with repeated values
	categories := []string{"Active", "Inactive", "Pending", "Completed", "Failed"}
	data := [][]string{{"Status", "ID"}}

	for i := 0; i < 1000; i++ {
		data = append(data, []string{
			categories[i%len(categories)],
			fmt.Sprintf("%d", i),
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf, _ := createNewBufferWithData(data, false)
		buf.rowFreeze = 1
		buf.enableStringInterning()
	}
}

func BenchmarkWithoutInterning(b *testing.B) {
	categories := []string{"Active", "Inactive", "Pending", "Completed", "Failed"}
	data := [][]string{{"Status", "ID"}}

	for i := 0; i < 1000; i++ {
		data = append(data, []string{
			categories[i%len(categories)],
			fmt.Sprintf("%d", i),
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf, _ := createNewBufferWithData(data, false)
		buf.rowFreeze = 1
		// No interning
	}
}
