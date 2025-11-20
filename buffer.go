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

// stringInterner provides efficient string deduplication for categorical data
type stringInterner struct {
	pool sync.Map // map[string]string for concurrent access
	hits uint64   // Cache hits (for debugging/stats)
}

// newStringInterner creates a new string interner
func newStringInterner() *stringInterner {
	return &stringInterner{}
}

// intern returns a canonical version of the string, reducing memory usage
func (si *stringInterner) intern(s string) string {
	// Quick path for empty strings
	if s == "" {
		return s
	}

	// Try to load existing string
	if existing, ok := si.pool.Load(s); ok {
		return existing.(string)
	}

	// Store and return the string
	si.pool.Store(s, s)
	return s
}

// shouldInternColumn determines if a column should use string interning
// based on its cardinality (ratio of unique values to total values)
func shouldInternColumn(values []string, threshold float64) bool {
	if len(values) < 100 {
		return false // Too small to benefit
	}

	// Sample the column to estimate cardinality
	sampleSize := 1000
	if len(values) < sampleSize {
		sampleSize = len(values)
	}

	seen := make(map[string]bool, sampleSize)
	for i := 0; i < sampleSize; i++ {
		seen[values[i]] = true
	}

	cardinality := float64(len(seen)) / float64(sampleSize)
	return cardinality < threshold // Low cardinality = good for interning
}

// Buffer represents a table data structure with concurrent access support
type Buffer struct {
	sep          rune              // Column separator character
	cont         [][]string        // Table content (rows x columns)
	colType      []int             // Column data types (colTypeStr or colTypeFloat)
	rowLen       int               // Number of rows
	colLen       int               // Number of columns
	rowFreeze    int               // Number of frozen header rows (0 or 1)
	colFreeze    int               // Number of frozen columns (0 or 1)
	selectedCell [][]int           // Selected cell coordinates
	mu           sync.RWMutex      // Mutex for concurrent access
	interners    []*stringInterner // String interners per column (nil if not used)
	internCols   []bool            // Track which columns use interning
	memoryUsage  int64             // Current estimated memory usage in bytes
	maxMemory    int64             // Maximum allowed memory in bytes (0 = no limit)
}

const (
	// Pre-allocated capacity for rows (optimized for large files)
	defaultRowCapacity = 10000
	// Cardinality threshold for string interning (30% unique values)
	internCardinalityThreshold = 0.30
	// Default memory limit: 0 = unlimited (users can set custom limit with --memory flag)
	defaultMaxMemoryBytes = 0
	// Estimated overhead per string in bytes (header + pointer + padding)
	stringOverheadBytes = 24
)

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
		interners:    nil,
		internCols:   nil,
		memoryUsage:  0,
		maxMemory:    defaultMaxMemoryBytes,
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

	// Check memory limit before adding row
	rowSize := b.estimateRowSize(s)
	if b.maxMemory > 0 && b.memoryUsage+rowSize > b.maxMemory {
		return errors.New("Memory limit exceeded: cannot load more data (limit: " +
			formatBytes(b.maxMemory) + ", current: " + formatBytes(b.memoryUsage) + ")")
	}

	// Strict mode: enforce column count
	if strict && len(s) != b.colLen {
		return errors.New("Row " + I2S(b.rowLen+b.rowFreeze) + " lacks some columns")
	}

	b.cont = append(b.cont, s)
	b.memoryUsage += rowSize

	// Adjust column count if needed
	if b.colLen != len(s) {
		b.resizeColUnsafe(len(s))
	}
	b.rowLen++

	return nil
}

// estimateRowSize estimates memory usage for a row in bytes
func (b *Buffer) estimateRowSize(row []string) int64 {
	size := int64(len(row) * 8) // Slice overhead (pointers)
	for _, s := range row {
		size += int64(len(s)) + stringOverheadBytes
	}
	return size
}

// getMemoryUsage returns current estimated memory usage
func (b *Buffer) getMemoryUsage() int64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.memoryUsage
}

// getMemoryLimit returns the configured memory limit
func (b *Buffer) getMemoryLimit() int64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.maxMemory
}

// setMemoryLimit sets the maximum memory limit in bytes (0 = no limit)
func (b *Buffer) setMemoryLimit(bytes int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.maxMemory = bytes
}

// getMemoryStats returns memory usage statistics
func (b *Buffer) getMemoryStats() map[string]interface{} {
	b.mu.RLock()
	defer b.mu.RUnlock()

	stats := make(map[string]interface{})
	stats["current_bytes"] = b.memoryUsage
	stats["current_formatted"] = formatBytes(b.memoryUsage)
	stats["limit_bytes"] = b.maxMemory
	stats["limit_formatted"] = formatBytes(b.maxMemory)

	if b.maxMemory > 0 {
		stats["usage_percent"] = float64(b.memoryUsage) / float64(b.maxMemory) * 100.0
		stats["available_bytes"] = b.maxMemory - b.memoryUsage
		stats["available_formatted"] = formatBytes(b.maxMemory - b.memoryUsage)
	} else {
		stats["usage_percent"] = 0.0
		stats["unlimited"] = true
	}

	return stats
}

// formatBytes formats bytes into human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return strconv.FormatInt(bytes, 10) + " B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"KB", "MB", "GB", "TB"}
	return strconv.FormatFloat(float64(bytes)/float64(div), 'f', 2, 64) + " " + units[exp]
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
		oldColLen := b.colLen
		b.colLen = n

		// Resize colType array if needed
		if len(b.colType) < n+1 {
			newColType := make([]int, n+1)
			copy(newColType, b.colType)
			b.colType = newColType
		}

		// Initialize new column types to colTypeStr (default)
		for i := oldColLen; i < n; i++ {
			b.colType[i] = colTypeStr
		}
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
// Returns 0 for invalid dates with fast pre-checks
func parseDateValueFast(s string) int64 {
	s = strings.TrimSpace(s)

	// Fast rejection checks
	if s == "" || s == "NA" || s == "N/A" || s == "null" {
		return 0
	}

	// Dates are typically 8-30 characters
	if len(s) < 8 || len(s) > 30 {
		return 0
	}

	// Must contain date separators
	if !strings.ContainsAny(s, "-/.:T ") {
		return 0
	}

	// Must contain at least one digit
	hasDigit := false
	for _, c := range s {
		if c >= '0' && c <= '9' {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return 0
	}

	// Try common date formats (most common first for performance)
	formats := []string{
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

	for _, format := range formats {
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

// isDateValue checks if a string represents a valid date with fast pre-checks
func isDateValue(s string) bool {
	if len(s) == 0 {
		return false
	}

	// Fast rejection: dates are typically 8-30 characters
	if len(s) < 8 || len(s) > 30 {
		return false
	}

	// Quick heuristic checks before trying to parse
	// Dates typically contain: -, /, :, T, or spaces with commas (for month names)
	hasDateSep := strings.ContainsAny(s, "-/.:T") || (strings.Contains(s, " ") && strings.Contains(s, ","))
	if !hasDateSep {
		return false
	}

	// Must contain at least one digit
	hasDigit := false
	for _, c := range s {
		if c >= '0' && c <= '9' {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return false
	}

	// Common date formats (most common first for performance)
	formats := []string{
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

	for _, format := range formats {
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

// detectAllColumnTypes automatically detects types for all columns in parallel
func (b *Buffer) detectAllColumnTypes() {
	types := make([]int, b.colLen)
	var wg sync.WaitGroup

	for i := 0; i < b.colLen; i++ {
		wg.Add(1)
		go func(col int) {
			defer wg.Done()
			types[col] = b.autoDetectColumnType(col)
		}(i)
	}

	wg.Wait()

	for i, t := range types {
		b.setColType(i, t)
	}
}

// enableStringInterning analyzes columns and enables interning for low-cardinality string columns
// This can save 30-70% memory for datasets with repeated categorical values
func (b *Buffer) enableStringInterning() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.rowLen < 100 {
		return // Too small to benefit
	}

	// Initialize interning structures
	b.interners = make([]*stringInterner, b.colLen)
	b.internCols = make([]bool, b.colLen)

	// Analyze each column
	for col := 0; col < b.colLen; col++ {
		// Skip non-string columns
		if b.colType[col] != colTypeStr {
			continue
		}

		// Get column data
		colData := make([]string, b.rowLen)
		for row := 0; row < b.rowLen; row++ {
			if col < len(b.cont[row]) {
				colData[row] = b.cont[row][col]
			}
		}

		// Check if column should be interned (low cardinality)
		if shouldInternColumn(colData, internCardinalityThreshold) {
			b.interners[col] = newStringInterner()
			b.internCols[col] = true

			// Intern existing values
			for row := 0; row < b.rowLen; row++ {
				if col < len(b.cont[row]) {
					b.cont[row][col] = b.interners[col].intern(b.cont[row][col])
				}
			}
		}
	}
}

// internValue interns a string value for a specific column if interning is enabled
func (b *Buffer) internValue(col int, value string) string {
	if col < len(b.internCols) && b.internCols[col] && b.interners[col] != nil {
		return b.interners[col].intern(value)
	}
	return value
}

// getInterningStats returns statistics about string interning usage
func (b *Buffer) getInterningStats() map[string]interface{} {
	b.mu.RLock()
	defer b.mu.RUnlock()

	stats := make(map[string]interface{})
	stats["enabled"] = len(b.internCols) > 0

	if len(b.internCols) == 0 {
		return stats
	}

	internedCols := 0
	for _, enabled := range b.internCols {
		if enabled {
			internedCols++
		}
	}

	stats["total_columns"] = b.colLen
	stats["interned_columns"] = internedCols
	stats["percentage"] = float64(internedCols) / float64(b.colLen) * 100.0

	return stats
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

	// Pre-allocate with estimated capacity (assume ~25% match rate)
	estimatedCapacity := (b.rowLen - b.rowFreeze) / 4
	if estimatedCapacity < 100 {
		estimatedCapacity = 100
	}
	filtered.cont = make([][]string, 0, estimatedCapacity)

	// Add header row if present
	if b.rowFreeze > 0 && b.rowLen > 0 {
		filtered.cont = append(filtered.cont, b.cont[0])
		filtered.rowLen = 1
	}

	// Early exit if column index is invalid - but still return buffer with header
	if colIndex >= b.colLen {
		return filtered
	}

	// Get column type for numeric comparisons
	colType := colTypeStr
	if colIndex < len(b.colType) {
		colType = b.colType[colIndex]
	}

	// Filter data rows
	startRow := b.rowFreeze
	for i := startRow; i < b.rowLen; i++ {
		if colIndex >= len(b.cont[i]) {
			continue
		}

		cellValue := b.cont[i][colIndex]

		// Evaluate filter condition
		if evaluateFilter(cellValue, options, colType) {
			filtered.cont = append(filtered.cont, b.cont[i])
			filtered.rowLen++
		}
	}

	return filtered
}

// evaluateFilter checks if a cell value matches the filter query based on the operator.
func evaluateFilter(cellValue string, options FilterOptions, colType int) bool {
	query := options.Query
	operator := options.Operator

	// Handle numeric comparisons first
	if colType == colTypeFloat || colType == colTypeDate {
		isNumericOperator := false
		switch operator {
		case ">", "<", ">=", "<=":
			isNumericOperator = true
		}

		if isNumericOperator {
			cellVal := parseNumericValueFast(cellValue)
			thresholdVal, err := strconv.ParseFloat(strings.TrimSpace(query), 64)
			if err != nil {
				return false // Cannot compare if query is not a number
			}

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
	}

	// Prepare strings for comparison
	cell := cellValue
	q := query
	if !options.CaseSensitive {
		cell = strings.ToLower(cell)
		q = strings.ToLower(q)
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
	case "regex":
		// When using regex, the user has full control over case sensitivity in the pattern.
		re, err := regexp.Compile(options.Query)
		if err != nil {
			return false // Invalid regex
		}
		return re.MatchString(cellValue)
	default:
		// Default to contains for backward compatibility if operator is empty
		return strings.Contains(cell, q)
	}
}
