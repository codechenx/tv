package main

import (
	"testing"
)

// ========================================
// Type Conversion Tests
// ========================================

func TestI2B(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Positive number 1", args{i: 1}, true},
		{"Positive number 2", args{i: 2}, true},
		{"Zero value", args{i: 0}, false},
		{"Negative value", args{i: -1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := I2B(tt.args.i); got != tt.want {
				t.Errorf("I2B() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestI2S_Basic(t *testing.T) {
	result := I2S(42)
	if result != "42" {
		t.Errorf("I2S(42) = %s, want 42", result)
	}
}

func TestF2S_Basic(t *testing.T) {
	result := F2S(42.0)
	if result == "" {
		t.Error("F2S should return a non-empty string")
	}
	t.Logf("F2S(42.0) = %s", result)
}

func TestS2F_Valid(t *testing.T) {
	result := S2F("3.14")
	if result != 3.14 {
		t.Errorf("S2F('3.14') = %f, want 3.14", result)
	}
}

func TestS2F_Invalid(t *testing.T) {
	t.Skip("Skipping S2F invalid test - function calls os.Exit on error")
}

// ========================================
// Text Wrapping Tests
// ========================================

func TestWrapText_Basic(t *testing.T) {
	result := wrapText("Hello World", 25)
	if len(result) < 1 {
		t.Error("wrapText should return at least one character")
	}
}

func TestWrapText_Long(t *testing.T) {
	longText := "This is a very long text that needs to be wrapped at twenty-five characters to test the wrapping functionality properly"
	result := wrapText(longText, 25)

	// Check that text was wrapped (contains newlines or is within limit)
	if len(longText) > 25 && len(result) == len(longText) {
		t.Error("Long text should be wrapped")
	}
}

func TestWrapText_Empty(t *testing.T) {
	result := wrapText("", 25)
	t.Logf("wrapText('', 25) returned: '%s'", result)
}

func TestWrapText_NoWrapNeeded(t *testing.T) {
	short := "Short"
	result := wrapText(short, 25)
	if result != short {
		t.Errorf("Short text should not be wrapped: got '%s', want '%s'", result, short)
	}
}

func TestWrapText_ExactLength(t *testing.T) {
	text := "Exactly25CharactersHere!!" // 25 characters
	result := wrapText(text, 25)
	if result != text {
		t.Errorf("Text at exact length should not be wrapped")
	}
}

// ========================================
// Text Truncation Tests
// ========================================

func TestTruncateText_Short(t *testing.T) {
	text := "Short"
	result := truncateText(text, 25)
	if result != text {
		t.Errorf("Short text should not be truncated: got '%s', want '%s'", result, text)
	}
}

func TestTruncateText_Long(t *testing.T) {
	text := "This is a very long text that needs to be truncated"
	result := truncateText(text, 20)
	expected := "This is a very lo..."
	if result != expected {
		t.Errorf("Long text should be truncated: got '%s', want '%s'", result, expected)
	}
	// Check that result doesn't exceed maxWidth
	if len([]rune(result)) > 20 {
		t.Errorf("Truncated text exceeds maxWidth: got length %d, want 20", len([]rune(result)))
	}
}

func TestTruncateText_ExactLength(t *testing.T) {
	text := "Exactly20Characters!"
	result := truncateText(text, 20)
	if result != text {
		t.Errorf("Text at exact length should not be truncated: got '%s', want '%s'", result, text)
	}
}

func TestTruncateText_ZeroWidth(t *testing.T) {
	text := "Some text"
	result := truncateText(text, 0)
	if result != text {
		t.Errorf("Zero maxWidth should return original text: got '%s', want '%s'", result, text)
	}
}

// ========================================
// Column Width Tests
// ========================================

func TestGetColumnMaxWidth_Valid(t *testing.T) {
	// Initialize wrapped columns map
	wrappedColumns = make(map[int]int)

	// Create a test buffer
	b = createNewBuffer()
	_ = b.contAppendSli([]string{"Short", "Medium", "VeryLongText"}, false)

	width := getColumnMaxWidth(0)
	if width < 1 {
		t.Error("getColumnMaxWidth should return positive width")
	}
}

func TestGetColumnMaxWidth_Custom(t *testing.T) {
	wrappedColumns = make(map[int]int)
	wrappedColumns[0] = 30

	width := getColumnMaxWidth(0)
	if width != 30 {
		t.Errorf("getColumnMaxWidth(0) = %d, want 30", width)
	}
}

func TestGetColumnMaxWidth_Default(t *testing.T) {
	wrappedColumns = make(map[int]int)

	width := getColumnMaxWidth(5)
	if width != 50 {
		t.Errorf("getColumnMaxWidth(5) = %d, want default 50", width)
	}
}

// ========================================
// Help Content Tests
// ========================================

func TestGetHelpContent_NotEmpty(t *testing.T) {
	help := getHelpContent()
	if len(help) == 0 {
		t.Error("getHelpContent() should return non-empty string")
	}
}

func TestGetHelpContent_ContainsBasics(t *testing.T) {
	help := getHelpContent()

	// Check for essential content
	if !contains(help, "Quit") {
		t.Error("Help should contain 'Quit' section")
	}
	if !contains(help, "Movement") {
		t.Error("Help should contain 'Movement' section")
	}
	if !contains(help, "Sort") {
		t.Error("Help should contain 'Sort' section")
	}
}

func TestUsefulInfo_NotEmpty(t *testing.T) {
	// Just verify it doesn't panic
	usefulInfo("test message")
	t.Log("usefulInfo executed successfully")
}

// ========================================
// Type Name Tests
// ========================================

func TestType2Name(t *testing.T) {
	tests := []struct {
		name     string
		colType  int
		expected string
	}{
		{"String type", colTypeStr, "Str"},
		{"Float type", colTypeFloat, "Num"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := type2name(tt.colType)
			if result != tt.expected {
				t.Errorf("type2name(%d) = %s, want %s", tt.colType, result, tt.expected)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
