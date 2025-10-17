package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//column data type
const colTypeStr = 0
const colTypeFloat = 1

//get column data type name. s: string, n: number
func type2name(i int) string {
	if i == 0 {
		return "Str"
	}
	return "Num"
}

var app *tview.Application
var UI *tview.Pages
var b *Buffer
var args Args
var debug bool
var statusMessage string       // Track status message for footer updates
var mainPage *tview.Frame      // Reference to main page for footer updates
var bufferTable *tview.Table   // Reference to buffer table
var fileNameStr string         // Store filename for footer
var cursorPosStr string        // Store cursor position for footer
var loadProgress LoadProgress  // Track loading progress
var userMovedCursor bool       // Track if user has moved the cursor

// LoadProgress tracks loading progress
type LoadProgress struct {
	TotalBytes   int64
	LoadedBytes  int64
	IsComplete   bool
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

// initialize tview, buffer
func initView() {
	app = tview.NewApplication()
	b = createNewBuffer()
}

//stop UI
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
