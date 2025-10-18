package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// column data type
const colTypeStr = 0
const colTypeFloat = 1
const colTypeDate = 2

// get column data type name. s: string, n: number, d: date
func type2name(i int) string {
	switch i {
	case colTypeStr:
		return "Str"
	case colTypeFloat:
		return "Num"
	case colTypeDate:
		return "Date"
	default:
		return "Str"
	}
}

var app *tview.Application
var UI *tview.Pages
var b *Buffer
var args Args
var debug bool
var statusMessage string         // Track status message for footer updates
var mainPage *tview.Frame        // Reference to main page for footer updates
var bufferTable *tview.Table     // Reference to buffer table
var fileNameStr string           // Store filename for footer
var cursorPosStr string          // Store cursor position for footer
var loadProgress LoadProgress    // Track loading progress
var userMovedCursor bool         // Track if user has moved the cursor
var wrappedColumns map[int]int   // Track which columns are wrapped and their max width
var searchResults []SearchResult // Store search results
var currentSearchIndex int       // Current position in search results
var searchQuery string           // Current search query
var searchModal tview.Primitive  // Search modal dialog
var searchUseRegex bool

var originalBuffer *Buffer              // Store original buffer before filtering
var isFiltered bool                     // Track if filter is active
var activeFilters map[int]FilterOptions // Track active filters: column -> query
var currentCursorColumn int             // Track current cursor column position
var lastKeyWasG bool                    // Track if last key pressed was 'g' for gg navigation

// LoadProgress tracks loading progress
type LoadProgress struct {
	TotalBytes  int64
	LoadedBytes int64
	IsComplete  bool
}

// GetPercentage returns the loading percentage (0-100)
func (lp *LoadProgress) GetPercentage() float64 {
	if lp.TotalBytes <= 0 {
		return 0
	}
	percent := float64(lp.LoadedBytes) * 100.0 / float64(lp.TotalBytes)
	if percent > 100 {
		percent = 100
	}
	return percent
}

// SearchResult represents a cell that matches search query
type SearchResult struct {
	Row int
	Col int
}

// initialize tview, buffer
func initView() {
	app = tview.NewApplication()
	app.EnableMouse(true) // Enable mouse support
	b = createNewBuffer()
	wrappedColumns = make(map[int]int) // Initialize wrapped columns map
	searchResults = []SearchResult{}
	currentSearchIndex = -1
	searchQuery = ""
	searchUseRegex = false
	originalBuffer = nil // Initialize filter variables
	isFiltered = false
	activeFilters = make(map[int]FilterOptions) // Initialize active filters map
	currentCursorColumn = 0                     // Initialize cursor column
	lastKeyWasG = false                         // Initialize vim navigation state
}

// stop UI
func stopView() {
	app.Stop()
}

// updateFooterWithStatus updates the footer with a status message
func updateFooterWithStatus(status string) {
	statusMessage = status
	if mainPage != nil {
		// Update the footer by rebuilding it
		mainPage.Clear()
		mainPage.AddText(fileNameStr, false, tview.AlignLeft, tcell.ColorDarkOrange).
			AddText(status, false, tview.AlignCenter, tcell.ColorDarkOrange).
			AddText(cursorPosStr, false, tview.AlignRight, tcell.ColorDarkOrange)
	}
}

//help page content
