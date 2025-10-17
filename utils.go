package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
)

// print fatal error and force quite app
func fatalError(err error) {
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println(err)
		color.Unset()
		if app != nil {
			app.Stop()
		}
		if !debug {
			os.Exit(1)
		}
	}
}

// print useful info and force quite app
func usefulInfo(s string) {
	color.Set(color.FgHiYellow)
	fmt.Println(s)
	color.Unset()
}

// I2B  covert int to bool, if i >0:true, else false
func I2B(i int) bool {
	return i > 0
}

// F2S covert float64 to bool
func F2S(i float64) string {
	return strconv.FormatFloat(i, 'f', 4, 64)
}

// S2F covert string to float64
func S2F(i string) float64 {
	s, err := strconv.ParseFloat(i, 64)
	if err != nil {
		fatalError(err)
	}
	return s
}

// I2S covert int to string
func I2S(i int) string {
	return strconv.Itoa(i)
}

func getHelpContent() string {
	helpContent := `
Vim-like Key Bindings

##Help Navigation##
j/k                 Scroll help text
gg/G                Jump to top/bottom
Ctrl-d/u            Page down/up
? or q or Esc       Close help dialog

##Quit##
q                   Quit application
Esc                 Close dialog or clear search

##Movement##
h                   Move left
l                   Move right
j                   Move down
k                   Move up

w                   Move to next column (word forward)
b                   Move to previous column (word backward)

gg                  Go to first row (press g twice)
G                   Go to last row

0                   Go to first column
$                   Go to last column

Ctrl-d              Page down (half page)
Ctrl-u              Page up (half page)

##Search##
/                   Search for text
n                   Next search result
N                   Previous search result
Esc                 Clear search highlighting

##Filter##
f                   Filter rows by current column value
r                   Reset/clear column filter

##Data Type##
t                   Toggle column data type (Str -> Num -> Date -> Str)

##Sort##
s                   Sort data by column (ascending)
S                   Sort data by column (descending)

##Text Wrapping##
W                   Toggle text wrapping for current column

##Stats##
i                   Show stats info for current column

##Help##
?                   Show this help dialog
`
	return helpContent
}

// wrapText wraps text to fit within maxWidth characters
// Returns the wrapped text with newlines
func wrapText(text string, maxWidth int) string {
	if maxWidth <= 0 || len(text) <= maxWidth {
		return text
	}

	var result []rune
	runes := []rune(text)
	lineStart := 0

	for i := 0; i < len(runes); i++ {
		// Check if we've reached the wrap point
		if i-lineStart >= maxWidth {
			// Find last space before maxWidth for word wrap
			wrapPoint := i
			for j := i; j > lineStart; j-- {
				if runes[j] == ' ' || runes[j] == '\t' || runes[j] == '-' {
					wrapPoint = j + 1
					break
				}
			}

			// If no good wrap point found, hard wrap at maxWidth
			if wrapPoint == i && i > lineStart {
				wrapPoint = lineStart + maxWidth
			}

			// Add the wrapped line
			result = append(result, runes[lineStart:wrapPoint]...)
			result = append(result, '\n')

			// Skip trailing spaces on new line
			for wrapPoint < len(runes) && (runes[wrapPoint] == ' ' || runes[wrapPoint] == '\t') {
				wrapPoint++
			}

			lineStart = wrapPoint
			i = wrapPoint - 1 // -1 because loop will increment
		}
	}

	// Add remaining text
	if lineStart < len(runes) {
		result = append(result, runes[lineStart:]...)
	}

	return string(result)
}

// getColumnMaxWidth determines the maximum width for a column
func getColumnMaxWidth(colIndex int) int {
	// Default wrap width (25 characters)
	defaultWidth := 25

	// Check if custom width is set
	if width, exists := wrappedColumns[colIndex]; exists {
		return width
	}

	return defaultWidth
}

// performSearch searches for a query string in the buffer and stores results
func performSearch(b *Buffer, query string, caseSensitive bool) []SearchResult {
	results := []SearchResult{}
	
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	searchQuery := query
	if !caseSensitive {
		searchQuery = toLower(query)
	}
	
	for r := 0; r < b.rowLen; r++ {
		for c := 0; c < b.colLen; c++ {
			cellText := b.cont[r][c]
			if !caseSensitive {
				cellText = toLower(cellText)
			}
			
			if stringContains(cellText, searchQuery) {
				results = append(results, SearchResult{Row: r, Col: c})
			}
		}
	}
	
	return results
}

// toLower converts a string to lowercase
func toLower(s string) string {
	runes := []rune(s)
	for i, r := range runes {
		if r >= 'A' && r <= 'Z' {
			runes[i] = r + 32
		}
	}
	return string(runes)
}

// stringContains checks if s contains substr
func stringContains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(substr) > len(s) {
		return false
	}
	
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// toLowerSimple converts a string to lowercase (simple implementation)
func toLowerSimple(s string) string {
	runes := []rune(s)
	for i, r := range runes {
		if r >= 'A' && r <= 'Z' {
			runes[i] = r + 32
		}
	}
	return string(runes)
}

// containsStr checks if s contains substr
func containsStr(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(substr) > len(s) {
		return false
	}
	
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
