package main

import (
	"testing"
)

func TestPerformSearch(t *testing.T) {
	// Create test buffer
	testData := [][]string{
		{"Name", "Age", "City"},
		{"John", "25", "New York"},
		{"Jane", "30", "Los Angeles"},
		{"Bob", "35", "Chicago"},
	}

	b, err := createNewBufferWithData(testData, false)
	if err != nil {
		t.Fatalf("Failed to create buffer: %v", err)
	}

	// Test case-insensitive search (non-regex)
	results := performSearch(b, "john", false)
	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'john', got %d", len(results))
	}
	if len(results) > 0 && (results[0].Row != 1 || results[0].Col != 0) {
		t.Errorf("Expected result at (1,0), got (%d,%d)", results[0].Row, results[0].Col)
	}

	// Test partial match (non-regex)
	results = performSearch(b, "an", false)
	if len(results) != 2 { // "Jane" and "Los Angeles"
		t.Errorf("Expected 2 results for 'an', got %d", len(results))
	}

	// Test no match (non-regex)
	results = performSearch(b, "xyz", false)
	if len(results) != 0 {
		t.Errorf("Expected 0 results for 'xyz', got %d", len(results))
	}
}

func TestPerformSearchRegex(t *testing.T) {
	// Create test buffer
	testData := [][]string{
		{"Name", "Email", "Age"},
		{"John", "john@test.com", "25"},
		{"Jane", "jane@example.org", "30"},
		{"Bob", "bob123@test.com", "35"},
	}

	b, err := createNewBufferWithData(testData, false)
	if err != nil {
		t.Fatalf("Failed to create buffer: %v", err)
	}

	// Test regex pattern for email domains
	results := performSearch(b, `@test\.com$`, true)
	if len(results) != 2 { // john@test.com and bob123@test.com
		t.Errorf("Expected 2 results for email regex, got %d", len(results))
	}

	// Test regex pattern for numbers at end
	results = performSearch(b, `\d+$`, true)
	if len(results) != 3 { // All age values (25, 30, 35)
		t.Errorf("Expected 3 results for number regex, got %d", len(results))
	}

	// Test regex pattern - case sensitive
	results = performSearch(b, `^J`, true)
	if len(results) != 2 { // John and Jane
		t.Errorf("Expected 2 results for '^J' regex, got %d", len(results))
	}

	// Test invalid regex
	results = performSearch(b, `[invalid(`, true)
	if len(results) != 0 {
		t.Errorf("Expected 0 results for invalid regex, got %d", len(results))
	}

	// Test regex OR pattern
	results = performSearch(b, `John|Bob`, true)
	if len(results) != 2 {
		t.Errorf("Expected 2 results for 'John|Bob' regex, got %d", len(results))
	}
}

func TestToLower(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HELLO", "hello"},
		{"Hello", "hello"},
		{"hello", "hello"},
		{"HeLLo WoRLd", "hello world"},
		{"123ABC", "123abc"},
		{"", ""},
	}

	for _, test := range tests {
		result := toLower(test.input)
		if result != test.expected {
			t.Errorf("toLower(%q) = %q, want %q", test.input, result, test.expected)
		}
	}
}

func TestStringContains(t *testing.T) {
	tests := []struct {
		s        string
		substr   string
		expected bool
	}{
		{"hello world", "world", true},
		{"hello world", "Hello", false},
		{"hello world", "o w", true},
		{"hello world", "", true},
		{"hello", "hello world", false},
		{"", "test", false},
		{"", "", true},
	}

	for _, test := range tests {
		result := stringContains(test.s, test.substr)
		if result != test.expected {
			t.Errorf("stringContains(%q, %q) = %v, want %v",
				test.s, test.substr, result, test.expected)
		}
	}
}
