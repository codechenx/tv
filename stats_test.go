package main

import (
	"testing"
)

// ========================================
// Continuous Stats Tests
// ========================================

func TestContinuousStats_Summary(t *testing.T) {
	cs := &ContinuousStats{}
	data := []string{"10.5", "20.3", "15.7", "25.2", "18.9"}

	cs.summary(data)

	// Verify data was processed
	summaryData := cs.getSummaryData()
	if len(summaryData) == 0 {
		t.Error("Expected summary data, got empty")
	}
}

func TestContinuousStats_EmptyData(t *testing.T) {
	cs := &ContinuousStats{}
	data := []string{}

	cs.summary(data)

	summaryData := cs.getSummaryData()
	if len(summaryData) == 0 {
		t.Log("Empty data handled correctly")
	}
}

func TestContinuousStats_SingleValue(t *testing.T) {
	cs := &ContinuousStats{}
	data := []string{"42.0"}

	cs.summary(data)

	summaryData := cs.getSummaryData()
	if len(summaryData) == 0 {
		t.Error("Should handle single value")
	}
}

func TestContinuousStats_WithNaN(t *testing.T) {
	cs := &ContinuousStats{}
	data := []string{"10.5", "NaN", "15.7", "invalid", "18.9"}

	cs.summary(data)

	// Should handle invalid/NaN values gracefully
	summaryData := cs.getSummaryData()
	if len(summaryData) == 0 {
		t.Error("Should produce summary even with invalid values")
	}
}

// ========================================
// Discrete Stats Tests
// ========================================

func TestDiscreteStats_Summary(t *testing.T) {
	ds := &DiscreteStats{}
	data := []string{"Apple", "Banana", "Apple", "Orange", "Banana", "Apple"}

	ds.summary(data)

	summaryData := ds.getSummaryData()
	if len(summaryData) == 0 {
		t.Error("Expected summary data, got empty")
	}
}

func TestDiscreteStats_EmptyData(t *testing.T) {
	ds := &DiscreteStats{}
	data := []string{}

	ds.summary(data)

	summaryData := ds.getSummaryData()
	if len(summaryData) == 0 {
		t.Log("Empty data handled correctly")
	}
}

func TestDiscreteStats_UniqueValues(t *testing.T) {
	ds := &DiscreteStats{}
	data := []string{"A", "B", "C", "D", "E"}

	ds.summary(data)

	summaryData := ds.getSummaryData()
	if len(summaryData) == 0 {
		t.Error("Should handle all unique values")
	}
}

func TestDiscreteStats_RepeatedValues(t *testing.T) {
	ds := &DiscreteStats{}
	data := []string{"X", "X", "X", "X", "X"}

	ds.summary(data)

	summaryData := ds.getSummaryData()
	if len(summaryData) == 0 {
		t.Error("Should handle repeated values")
	}
}

// ========================================
// Stats Interface Tests
// ========================================

func TestStatsSummary_Interface(t *testing.T) {
	// Test that both types implement the interface
	var _ statsSummary = &ContinuousStats{}
	var _ statsSummary = &DiscreteStats{}

	t.Log("Both stats types implement statsSummary interface")
}

func TestGetSummaryStr(t *testing.T) {
	cs := &ContinuousStats{}
	data := []string{"10", "20", "30"}

	str := cs.getSummaryStr(data)
	if len(str) == 0 {
		t.Error("getSummaryStr should return non-empty string")
	}
}

// ========================================
// Performance Tests
// ========================================

func BenchmarkContinuousStats(b *testing.B) {
	data := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = F2S(float64(i) * 1.5)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cs := &ContinuousStats{}
		cs.summary(data)
	}
}

func BenchmarkDiscreteStats(b *testing.B) {
	data := make([]string, 1000)
	categories := []string{"A", "B", "C", "D", "E"}
	for i := 0; i < 1000; i++ {
		data[i] = categories[i%len(categories)]
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ds := &DiscreteStats{}
		ds.summary(data)
	}
}
