package main

import (
	"errors"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Buffer represents a table data structure with concurrent access support
type Buffer struct {
	sep          rune         // Column separator character
	cont         [][]string   // Table content (rows x columns)
	colType      []int        // Column data types (colTypeStr or colTypeFloat)
	rowLen       int          // Number of rows
	colLen       int          // Number of columns
	rowFreeze    int          // Number of frozen header rows (0 or 1)
	colFreeze    int          // Number of frozen columns (0 or 1)
	selectedCell [][]int      // Selected cell coordinates
	mu           sync.RWMutex // Mutex for concurrent access
}

const (
	// Pre-allocated capacity for rows (optimized for large files)
	defaultRowCapacity = 10000
)

// Common date formats for parsing (most common first for performance)
// Shared by both isDateValue and parseDateValueFast to ensure consistency
var commonDateFormats = []string{
	"2006-01-02",          // ISO date: 2024-10-17
	"2006-01-02 15:04:05", // ISO datetime: 2024-10-17 15:30:00
	"01/02/2006",          // US date: 10/17/2024
	"02/01/2006",          // EU date: 17/10/2024
	"2006/01/02",          // Alt ISO: 2024/10/17
	time.RFC3339,          // RFC3339: 2024-10-17T15:30:00Z
	"2006-01-02T15:04:05", // ISO8601 without timezone
	"Jan 02, 2006",        // Mon DD, YYYY
	"January 02, 2006",    // Month DD, YYYY
	"02-Jan-2006",         // DD-Mon-YYYY
	"02 Jan 2006",         // DD Mon YYYY
	"2006.01.02",          // Dotted date
}

// createNewBuffer initializes and returns a new empty Buffer
func createNewBuffer() *Buffer {
	return &Buffer{
		sep:          0,
		cont:         [][]string{},
		colType:      []int{},
		rowLen:       0,
		colLen:       0,
		rowFreeze:    1,
		colFreeze:    1,
		selectedCell: [][]int{},
	}
}

// createNewBufferWithData creates a Buffer from existing data
func createNewBufferWithData(ss [][]string, strict bool) (*Buffer, error) {
	b = createNewBuffer()
	for _, s := range ss {
		if err := b.contAppendSli(s, strict); err != nil {
			return nil, err
		}
	}
	return b, nil
}

// contAppendSli appends a row to the buffer
// strict: if true, enforces consistent column count
func (b *Buffer) contAppendSli(s []string, strict bool) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Initialize on first row
	if b.rowLen == 0 {
		b.colLen = len(s)
		b.colType = make([]int, b.colLen+1)
		// Pre-allocate capacity to reduce reallocations
		if cap(b.cont) == 0 {
			b.cont = make([][]string, 0, defaultRowCapacity)
		}
	}

	// Strict mode: enforce column count
	if strict && len(s) != b.colLen {
		return errors.New("Row " + I2S(b.rowLen+b.rowFreeze) + " lacks some columns")
	}

	b.cont = append(b.cont, s)

	// Adjust column count if needed
	if b.colLen != len(s) {
		b.resizeColUnsafe(len(s))
	}
	b.rowLen++

	return nil
}

// resizeColUnsafe adjusts the number of columns (must be called with lock held)
// Fills missing columns with "NaN"
func (b *Buffer) resizeColUnsafe(n int) {
	if n <= 0 {
		return
	}

	lackLen := b.colLen - n
	if lackLen < 0 {
		lackLen = n - b.colLen
		b.colLen = n
	}

	// Fill missing columns with NaN
	for ii := range b.cont {
		for m := 0; m < lackLen; m++ {
			b.cont[ii] = append(b.cont[ii], "NaN")
		}
	}
}

// resizeCol adjusts the number of columns (thread-safe)
func (b *Buffer) resizeCol(n int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.resizeColUnsafe(n)
}

// sortByStr sorts the buffer by column in string mode
// colIndex: column to sort by
// rev: true for descending, false for ascending
func (b *Buffer) sortByStr(colIndex int, rev bool) {
	hasHeader := I2B(b.rowFreeze)

	if rev {
		// Descending sort
		if hasHeader {
			sort.SliceStable(b.cont[1:], func(i, j int) bool {
				return b.cont[1:][i][colIndex] > b.cont[1:][j][colIndex]
			})
		} else {
			sort.SliceStable(b.cont, func(i, j int) bool {
				return b.cont[i][colIndex] > b.cont[j][colIndex]
			})
		}
	} else {
		// Ascending sort
		if hasHeader {
			sort.SliceStable(b.cont[1:], func(i, j int) bool {
				return b.cont[1:][i][colIndex] < b.cont[1:][j][colIndex]
			})
		} else {
			sort.SliceStable(b.cont, func(i, j int) bool {
				return b.cont[i][colIndex] < b.cont[j][colIndex]
			})
		}
	}
}

// sortByNum sorts column by number format with optimized numeric conversion
func (b *Buffer) sortByNum(colIndex int, rev bool) {
	hasHeader := I2B(b.rowFreeze)
	dataRows := b.cont
	if hasHeader {
		dataRows = b.cont[1:]
	}

	// Create index-value pairs to sort
	type numRow struct {
		row []string
		num float64
	}

	pairs := make([]numRow, len(dataRows))
	for i := range dataRows {
		pairs[i] = numRow{
			row: dataRows[i],
			num: parseNumericValueFast(dataRows[i][colIndex]),
		}
	}

	// Sort the pairs
	if rev {
		sort.SliceStable(pairs, func(i, j int) bool {
			return pairs[i].num > pairs[j].num
		})
	} else {
		sort.SliceStable(pairs, func(i, j int) bool {
			return pairs[i].num < pairs[j].num
		})
	}

	// Copy back sorted rows
	for i := range pairs {
		dataRows[i] = pairs[i].row
	}
}

// sortByDate sorts column by date format with optimized date parsing
func (b *Buffer) sortByDate(colIndex int, rev bool) {
	hasHeader := I2B(b.rowFreeze)
	dataRows := b.cont
	if hasHeader {
		dataRows = b.cont[1:]
	}

	// Create index-value pairs to sort
	type dateRow struct {
		row  []string
		date int64
	}

	pairs := make([]dateRow, len(dataRows))
	for i := range dataRows {
		pairs[i] = dateRow{
			row:  dataRows[i],
			date: parseDateValueFast(dataRows[i][colIndex]),
		}
	}

	// Sort the pairs
	if rev {
		sort.SliceStable(pairs, func(i, j int) bool {
			return pairs[i].date > pairs[j].date
		})
	} else {
		sort.SliceStable(pairs, func(i, j int) bool {
			return pairs[i].date < pairs[j].date
		})
	}

	// Copy back sorted rows
	for i := range pairs {
		dataRows[i] = pairs[i].row
	}
}

// parseNumericValueFast quickly parses a string to float64
// Handles commas, underscores, and returns 0 for invalid values
func parseNumericValueFast(s string) float64 {
	// Remove common separators
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.TrimSpace(s)

	if s == "" || s == "NA" || s == "N/A" || s == "NaN" || s == "null" {
		return 0
	}

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return val
}

// parseDateValueFast quickly parses a date string to unix timestamp
// Returns 0 for invalid dates
func parseDateValueFast(s string) int64 {
	s = strings.TrimSpace(s)

	if s == "" || s == "NA" || s == "N/A" || s == "null" {
		return 0
	}

	// Try common date formats using shared constant
	for _, format := range commonDateFormats {
		if t, err := time.Parse(format, s); err == nil {
			return t.Unix()
		}
	}

	return 0
}

// getCol returns the ith column data as a string slice
// Uses pointer receiver to avoid copying mutex
func (b *Buffer) getCol(i int) []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result := make([]string, b.rowLen)
	for rowI := 0; rowI < b.rowLen; rowI++ {
		result[rowI] = b.cont[rowI][i]
	}
	return result
}

// set ith column data type
func (b *Buffer) setColType(i int, t int) {
	b.colType[i] = t
}

// get ith column data type
func (b *Buffer) getColType(i int) int {
	return b.colType[i]
}

// autoDetectColumnType intelligently detects if a column contains numeric, date, or string data
// Returns colTypeDate for dates, colTypeFloat for numbers, colTypeStr for strings
func (b *Buffer) autoDetectColumnType(colIndex int) int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if colIndex < 0 || colIndex >= b.colLen {
		return colTypeStr
	}

	// Sample size for type detection
	startRow := b.rowFreeze
	endRow := b.rowLen

	// For large datasets, sample smartly (first N rows + some middle + last N)
	sampleSize := 100
	sampleRows := []int{}

	if endRow-startRow > sampleSize {
		// Sample first 50 rows
		for i := startRow; i < startRow+50 && i < endRow; i++ {
			sampleRows = append(sampleRows, i)
		}
		// Sample middle 25 rows
		midPoint := (startRow + endRow) / 2
		for i := midPoint; i < midPoint+25 && i < endRow; i++ {
			sampleRows = append(sampleRows, i)
		}
		// Sample last 25 rows
		for i := endRow - 25; i < endRow; i++ {
			if i > startRow {
				sampleRows = append(sampleRows, i)
			}
		}
	} else {
		// For small datasets, check all rows
		for i := startRow; i < endRow; i++ {
			sampleRows = append(sampleRows, i)
		}
	}

	// Analyze samples
	dateCount := 0
	numericCount := 0
	totalCount := 0

	for _, rowIdx := range sampleRows {
		if rowIdx >= b.rowLen || colIndex >= len(b.cont[rowIdx]) {
			continue
		}

		value := strings.TrimSpace(b.cont[rowIdx][colIndex])

		// Skip empty/null cells
		if value == "" || value == "NA" || value == "N/A" || value == "NaN" || value == "null" {
			continue
		}

		totalCount++

		// Check if it's a date (dates are more specific than numbers)
		if isDateValue(value) {
			dateCount++
		} else if isNumericValue(value) {
			numericCount++
		}
	}

	// If no valid values, treat as string
	if totalCount == 0 {
		return colTypeStr
	}

	// Threshold: 90% of values must match type
	threshold := float64(totalCount) * 0.90

	// Priority: Date > Number > String
	if float64(dateCount) >= threshold {
		return colTypeDate
	} else if float64(numericCount) >= threshold {
		return colTypeFloat
	}

	return colTypeStr
}

// isDateValue checks if a string represents a valid date
func isDateValue(s string) bool {
	if len(s) == 0 {
		return false
	}

	// Quick heuristic checks before trying to parse
	// Dates typically contain: -, /, :, T, or spaces with commas (for month names)
	hasDateSep := strings.ContainsAny(s, "-/.:T") || (strings.Contains(s, " ") && strings.Contains(s, ","))
	if !hasDateSep {
		return false
	}

	// Try common date formats using shared constant
	for _, format := range commonDateFormats {
		if _, err := time.Parse(format, s); err == nil {
			return true
		}
	}

	return false
}

// isNumericValue checks if a string represents a valid number
// Handles: integers, floats, scientific notation, negative numbers
func isNumericValue(s string) bool {
	if len(s) == 0 {
		return false
	}

	// Quick check for common patterns
	hasDigit := false
	hasDot := false
	hasE := false
	i := 0

	// Handle sign
	if s[i] == '+' || s[i] == '-' {
		i++
		if i >= len(s) {
			return false
		}
	}

	// Parse number
	for i < len(s) {
		c := s[i]

		if c >= '0' && c <= '9' {
			hasDigit = true
		} else if c == '.' {
			if hasDot || hasE {
				return false // Multiple dots or dot after E
			}
			hasDot = true
		} else if c == 'e' || c == 'E' {
			if !hasDigit || hasE {
				return false // E without digits or multiple E
			}
			hasE = true
			hasDigit = false // Reset for exponent part

			// Check for sign after E
			if i+1 < len(s) && (s[i+1] == '+' || s[i+1] == '-') {
				i++
			}
		} else if c == '_' || c == ',' {
			// Allow thousand separators (common in data files)
			// Skip validation, just continue
		} else {
			return false // Invalid character
		}
		i++
	}

	return hasDigit
}

// detectAllColumnTypes automatically detects types for all columns
// Uses parallel processing for better performance on multi-column datasets
func (b *Buffer) detectAllColumnTypes() {
	// For small number of columns, sequential processing is faster
	if b.colLen <= 4 {
		for i := 0; i < b.colLen; i++ {
			detectedType := b.autoDetectColumnType(i)
			b.setColType(i, detectedType)
		}
		return
	}

	// For larger datasets, use parallel processing
	type result struct {
		index int
		ctype int
	}

	results := make(chan result, b.colLen)
	var wg sync.WaitGroup

	// Process columns in parallel
	for i := 0; i < b.colLen; i++ {
		wg.Add(1)
		go func(colIndex int) {
			defer wg.Done()
			detectedType := b.autoDetectColumnType(colIndex)
			results <- result{index: colIndex, ctype: detectedType}
		}(i)
	}

	// Close results channel when all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results and set column types
	for res := range results {
		b.setColType(res.index, res.ctype)
	}
}

//clear selectedCell of buffer
//func (b *Buffer) clearSelection() {
//	b.selectedCell = [][]int{}
//}

// search string and add result to selectedCell of buffer
func (b *Buffer) selectBySearch(s string) {
	for ii, i := range b.cont {
		for ji, j := range i {
			if s == j {
				b.selectedCell = append(b.selectedCell, []int{ii, ji})
			}
		}
	}
}

// FilterOptions defines the parameters for a column filter.
type FilterOptions struct {
	Query         string
	Operator      string
	CaseSensitive bool
}

// filterByColumn filters rows based on column value using the provided options.
// It returns a new buffer containing the filtered rows.
func (b *Buffer) filterByColumn(colIndex int, options FilterOptions) *Buffer {
	b.mu.RLock()
	defer b.mu.RUnlock()

	filtered := createNewBuffer()
	filtered.sep = b.sep
	filtered.colLen = b.colLen
	filtered.rowFreeze = b.rowFreeze
	filtered.colFreeze = b.colFreeze
	filtered.colType = make([]int, len(b.colType))
	copy(filtered.colType, b.colType)

	// Add header row if present
	if b.rowFreeze > 0 && b.rowLen > 0 {
		filtered.cont = append(filtered.cont, b.cont[0])
		filtered.rowLen = 1
	}

	// Get column type for numeric comparisons
	colType := colTypeStr
	if colIndex < len(b.colType) {
		colType = b.colType[colIndex]
	}

	// Pre-compile regex if using regex operator (performance optimization)
	var compiledRegex *regexp.Regexp
	if options.Operator == "regex" {
		var err error
		compiledRegex, err = regexp.Compile(options.Query)
		if err != nil {
			// Invalid regex - return buffer with just header
			return filtered
		}
	}

	// Pre-parse numeric threshold if using numeric operators (performance optimization)
	var thresholdVal float64
	isNumericOp := false
	if colType == colTypeFloat || colType == colTypeDate {
		switch options.Operator {
		case ">", "<", ">=", "<=":
			isNumericOp = true
			var err error
			thresholdVal, err = strconv.ParseFloat(strings.TrimSpace(options.Query), 64)
			if err != nil {
				// Invalid numeric threshold - return buffer with just header
				return filtered
			}
		}
	}

	// Pre-convert query to lowercase if case-insensitive (performance optimization)
	lowerQuery := options.Query
	if !options.CaseSensitive && options.Operator != "regex" {
		lowerQuery = strings.ToLower(options.Query)
	}

	// Filter data rows
	startRow := b.rowFreeze
	for i := startRow; i < b.rowLen; i++ {
		if colIndex >= len(b.cont[i]) {
			continue
		}

		cellValue := b.cont[i][colIndex]

		// Evaluate filter condition with pre-compiled/parsed values
		if evaluateFilterOptimized(cellValue, options, colType, compiledRegex, isNumericOp, thresholdVal, lowerQuery) {
			filtered.cont = append(filtered.cont, b.cont[i])
			filtered.rowLen++
		}
	}

	return filtered
}

// evaluateFilterOptimized checks if a cell value matches the filter query based on the operator.
// This version accepts pre-compiled regex and pre-parsed values for better performance.
func evaluateFilterOptimized(cellValue string, options FilterOptions, colType int, compiledRegex *regexp.Regexp, isNumericOp bool, thresholdVal float64, lowerQuery string) bool {
	operator := options.Operator

	// Handle numeric comparisons first
	if isNumericOp {
		cellVal := parseNumericValueFast(cellValue)
		switch operator {
		case ">":
			return cellVal > thresholdVal
		case "<":
			return cellVal < thresholdVal
		case ">=":
			return cellVal >= thresholdVal
		case "<=":
			return cellVal <= thresholdVal
		}
	}

	// Handle regex operator with pre-compiled regex
	if operator == "regex" && compiledRegex != nil {
		return compiledRegex.MatchString(cellValue)
	}

	// Prepare strings for comparison
	cell := cellValue
	q := lowerQuery
	if !options.CaseSensitive {
		cell = strings.ToLower(cell)
	}

	// Handle string-based operators
	switch operator {
	case "contains":
		return strings.Contains(cell, q)
	case "equals":
		return cell == q
	case "starts with":
		return strings.HasPrefix(cell, q)
	case "ends with":
		return strings.HasSuffix(cell, q)
	default:
		// Default to contains for backward compatibility if operator is empty
		return strings.Contains(cell, q)
	}
}

// evaluateFilter checks if a cell value matches the filter query based on the operator.
// Kept for backward compatibility - calls evaluateFilterOptimized with default parameters.
func evaluateFilter(cellValue string, options FilterOptions, colType int) bool {
	// For backward compatibility, compile regex on-the-fly if needed
	var compiledRegex *regexp.Regexp
	if options.Operator == "regex" {
		var err error
		compiledRegex, err = regexp.Compile(options.Query)
		if err != nil {
			return false
		}
	}

	// Parse numeric threshold if needed
	var thresholdVal float64
	isNumericOp := false
	if colType == colTypeFloat || colType == colTypeDate {
		switch options.Operator {
		case ">", "<", ">=", "<=":
			isNumericOp = true
			var err error
			thresholdVal, err = strconv.ParseFloat(strings.TrimSpace(options.Query), 64)
			if err != nil {
				return false
			}
		}
	}

	// Convert query to lowercase if needed
	lowerQuery := options.Query
	if !options.CaseSensitive && options.Operator != "regex" {
		lowerQuery = strings.ToLower(options.Query)
	}

	return evaluateFilterOptimized(cellValue, options, colType, compiledRegex, isNumericOp, thresholdVal, lowerQuery)
}
