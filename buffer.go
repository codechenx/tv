package main

import (
	"errors"
	"sort"
	"sync"
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

// sort column by number format
func (b *Buffer) sortByNum(colIndex int, rev bool) {
	if rev {
		if I2B(b.rowFreeze) {
			sort.SliceStable(b.cont[1:], func(i, j int) bool { return S2F(b.cont[1:][i][colIndex]) > S2F(b.cont[1:][j][colIndex]) })
		} else {
			sort.SliceStable(b.cont, func(i, j int) bool { return S2F(b.cont[i][colIndex]) > S2F(b.cont[j][colIndex]) })
		}
	} else {

		if I2B(b.rowFreeze) {
			sort.SliceStable(b.cont[1:], func(i, j int) bool { return S2F(b.cont[1:][i][colIndex]) < S2F(b.cont[1:][j][colIndex]) })
		} else {
			sort.SliceStable(b.cont, func(i, j int) bool { return S2F(b.cont[i][colIndex]) < S2F(b.cont[j][colIndex]) })
		}
	}
}

// transpose buffer content
func (b *Buffer) transpose() {
	b.rowLen, b.colLen = b.colLen, b.rowLen
	b.colType = make([]int, b.colLen+1)
	xl := len(b.cont[0])
	yl := len(b.cont)
	result := make([][]string, xl)
	for i := range result {
		result[i] = make([]string, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = b.cont[j][i]
		}
	}
	b.cont = result
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

// filterByColumn filters rows based on column value containing the search string
// Returns a new filtered buffer
func (b *Buffer) filterByColumn(colIndex int, query string, caseSensitive bool) *Buffer {
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
	if b.rowFreeze > 0 {
		filtered.cont = append(filtered.cont, b.cont[0])
		filtered.rowLen = 1
	}

	// Filter data rows
	startRow := b.rowFreeze
	for i := startRow; i < b.rowLen; i++ {
		if colIndex >= len(b.cont[i]) {
			continue
		}
		
		cellValue := b.cont[i][colIndex]
		queryStr := query
		
		if !caseSensitive {
			cellValue = toLowerSimple(cellValue)
			queryStr = toLowerSimple(query)
		}
		
		if containsStr(cellValue, queryStr) {
			filtered.cont = append(filtered.cont, b.cont[i])
			filtered.rowLen++
		}
	}

	return filtered
}
