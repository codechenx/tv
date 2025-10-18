package main

import (
	"fmt"
	"os"
	"regexp"
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
	helpContent := `[::b][yellow]━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━[white]

[::b][cyan]🚀 TV - Modern Terminal Table Viewer[-][white]

[::b][yellow]━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━[white]

[::b][green]📖 Help Navigation[white]
  [yellow]j/k[-]                 Scroll help text
  [yellow]gg/G[-]                Jump to top/bottom
  [yellow]Ctrl-d/u[-]            Page down/up
  [yellow]? or q or Esc[-]       Close help dialog

[::b][red]🚪 Quit[white]
  [yellow]q[-]                   Quit application
  [yellow]Esc[-]                 Close dialog or clear search

[::b][blue]⬆️ Movement[white]
  [yellow]h[-]                   Move left ⬅️
  [yellow]l[-]                   Move right ➡️
  [yellow]j[-]                   Move down ⬇️
  [yellow]k[-]                   Move up ⬆️

  [yellow]w[-]                   Move to next column (word forward)
  [yellow]b[-]                   Move to previous column (word backward)

  [yellow]gg[-]                  Go to first row (press g twice)
  [yellow]G[-]                   Go to last row

  [yellow]0[-]                   Go to first column
  [yellow]$[-]                   Go to last column

  [yellow]Ctrl-d[-]              Page down (half page)
  [yellow]Ctrl-u[-]              Page up (half page)

[::b][magenta]🔍 Search[white]
  [yellow]/[-]                   Search for text
                    • Case-insensitive by default
                    • Press [yellow]Tab[-] to navigate to checkbox
                    • Press [yellow]Space[-] to toggle [yellow]Use Regex[-] option
  [yellow]n[-]                   Next search result ⏭
  [yellow]N[-]                   Previous search result ⏮
  [yellow]Esc[-]                 Clear search highlighting

[::b][green]🎯 Regex Search Examples[white]
  [yellow]^start[-]              Match at beginning of cell
  [yellow]end$[-]                Match at end of cell
  [yellow]\d+[-]                 Match digits (numbers)
  [yellow]@.*\.com[-]            Match email pattern
  [yellow]word1|word2[-]         Match either word (OR)
  [yellow][A-Z]+[-]              Match uppercase letters

[::b][orange]🔎 Filter[white]
  [yellow]f[-]                   Filter rows by current column value
                    • Apply filters to multiple columns
                    • Edit filter: press f on filtered column
                    OR: same cell has either term
                    AND: same cell has both terms
                    ROR: different rows, any match (uppercase only)
  [yellow]r[-]                   Remove filter from current column

[::b][purple]🏷️  Data Type[white]
  [yellow]t[-]                   Toggle column data type
                    (String → Number → Date → String)

[::b][green]🔃 Sort[white]
  [yellow]s[-]                   Sort data by column (ascending ⬆️)
  [yellow]S[-]                   Sort data by column (descending ⬇️)

[::b][cyan]📏 Text Wrapping[white]
  [yellow]W[-]                   Toggle width limit for current column (50 chars)
                    Long columns (>50 chars) are limited automatically

[::b][blue]📊 Stats[white]
  [yellow]i[-]                   Show stats info for current column

[::b][yellow]❓ Help[white]
  [yellow]?[-]                   Show this help dialog

[::b][yellow]━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━[white]

[::b][green]💡 Pro Tips:[white]
  • Press [yellow]gg[-] to jump to the top of any table
  • Use [yellow]/[-] for quick searching across all cells
  • Enable [yellow]regex[-] mode for powerful pattern matching
  • Press [yellow]i[-] to see detailed statistics for any column
  • Use [yellow]f[-] on multiple columns to combine filters
  • Headers are frozen by default for easy navigation

[::b][yellow]━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━[white]
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

// truncateText truncates text to maxWidth and adds ellipsis if needed
func truncateText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}

	runes := []rune(text)
	if len(runes) <= maxWidth {
		return text
	}

	// Reserve 3 characters for ellipsis
	if maxWidth <= 3 {
		return string(runes[:maxWidth])
	}

	return string(runes[:maxWidth-3]) + "..."
}

// getColumnMaxWidth determines the maximum width for a column
func getColumnMaxWidth(colIndex int) int {
	// Default wrap width (50 characters for long columns)
	defaultWidth := 50

	// Check if custom width is set
	if width, exists := wrappedColumns[colIndex]; exists {
		return width
	}

	return defaultWidth
}

// detectAndWrapLongColumns automatically enables wrapping for columns with long content
// Analyzes first N rows to detect if columns have text longer than threshold
func detectAndWrapLongColumns(b *Buffer, sampleSize int, threshold int) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// Determine how many rows to sample
	maxSample := sampleSize
	if b.rowLen < maxSample {
		maxSample = b.rowLen
	}

	// Skip header row in analysis if it exists
	startRow := 0
	if b.rowFreeze > 0 {
		startRow = b.rowFreeze
	}

	// Track maximum length found in each column
	maxLengths := make([]int, b.colLen)

	// Sample rows to find maximum content length per column
	for r := startRow; r < maxSample; r++ {
		for c := 0; c < b.colLen; c++ {
			if c < len(b.cont[r]) {
				cellLen := len(b.cont[r][c])
				if cellLen > maxLengths[c] {
					maxLengths[c] = cellLen
				}
			}
		}
	}

	// Enable wrapping for columns that exceed threshold
	for c := 0; c < b.colLen; c++ {
		if maxLengths[c] > threshold {
			// Only set if not already manually configured
			if _, exists := wrappedColumns[c]; !exists {
				wrappedColumns[c] = getColumnMaxWidth(c)
			}
		}
	}
}

// performSearch searches for a query string in the buffer and stores results
// Supports both plain text and regex search modes
func performSearch(b *Buffer, query string, useRegex bool) []SearchResult {
	results := []SearchResult{}

	b.mu.RLock()
	defer b.mu.RUnlock()

	// Compile regex if in regex mode
	var re *regexp.Regexp
	var err error
	if useRegex {
		re, err = regexp.Compile(query)
		if err != nil {
			// If regex is invalid, return empty results
			return results
		}
	} else {
		// For non-regex, convert to lowercase for case-insensitive search
		query = toLower(query)
	}

	// Scan column by column (same column first, then next column)
	for c := 0; c < b.colLen; c++ {
		for r := 0; r < b.rowLen; r++ {
			cellText := b.cont[r][c]
			
			var matches bool
			if useRegex {
				matches = re.MatchString(cellText)
			} else {
				matches = stringContains(toLower(cellText), query)
			}

			if matches {
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
