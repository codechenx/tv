package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// buildCursorPosStr builds the cursor position string (without filter info now)
func buildCursorPosStr(row, column int) string {
	posStr := "Column Type: " + type2name(b.getColType(column)) + "  |  " + strconv.Itoa(row) + "," + strconv.Itoa(column) + "  "
	return posStr
}

// buildFilterInfoStr builds the filter information string for the top strip
// Shows all active filters or current column filter when cursor is on a filtered column
func buildFilterInfoStr(currentColumn int) string {
	if !isFiltered || len(activeFilters) == 0 {
		return "" // No filter active
	}

	// Check if current column has a filter
	if opts, hasFilter := activeFilters[currentColumn]; hasFilter {
		// Get column name if available
		columnName := fmt.Sprintf("Column %d", currentColumn)
		if b.rowFreeze > 0 && len(b.cont) > 0 && currentColumn < len(b.cont[0]) {
			columnName = b.cont[0][currentColumn]
		}

		return fmt.Sprintf("🔎 Filter Active: [%s] %s \"%s\"  |  %d filters total  |  Press 'r' to remove this filter", columnName, opts.Operator, opts.Query, len(activeFilters))
	}

	// Show summary if cursor is not on a filtered column
	return fmt.Sprintf("🔎 %d filters active  |  Navigate to filtered column and press 'r' to remove", len(activeFilters))
}

// add buffer data to buffer table
func drawBuffer(b *Buffer, t *tview.Table) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	t.Clear()
	cols, rows := b.colLen, b.rowLen

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			backgroundColor := tcell.ColorDefault
			attributes := tcell.AttrNone
			alignment := tview.AlignLeft

			// Get cell content
			cellText := b.cont[r][c]

			// Check if this is a header row/column (frozen area)
			isHeaderRow := r < b.rowFreeze && args.Header != -1 && args.Header != 2
			isHeaderCol := c < b.colFreeze

			// Modern header styling with rich visual design
			if isHeaderRow {
				// Main header row: bold white text on gradient blue background
				color = tcell.ColorWhite
				backgroundColor = tcell.NewRGBColor(30, 60, 120) // Deep blue
				attributes = tcell.AttrBold | tcell.AttrUnderline
				alignment = tview.AlignCenter

				// Add filter indicator if this column has a filter applied
				if isFiltered {
					if _, hasFilter := activeFilters[c]; hasFilter {
						cellText = "🔎 " + cellText + " 🔎"
						backgroundColor = tcell.NewRGBColor(255, 100, 0) // Orange background for filtered column
					}
				}
			} else if isHeaderCol {
				// Frozen column: gold color for row headers
				color = tcell.NewRGBColor(255, 215, 0) // Gold
				attributes = tcell.AttrBold
			}

			// Check if this cell is a search result and highlight it
			isSearchMatch := false
			if searchQuery != "" && len(searchResults) > 0 {
				for _, result := range searchResults {
					if result.Row == r && result.Col == c {
						isSearchMatch = true
						break
					}
				}
			}

			// Modern search match highlighting (overrides header styling)
			if isSearchMatch {
				// Check if this is the current search result
				if currentSearchIndex >= 0 && currentSearchIndex < len(searchResults) &&
					searchResults[currentSearchIndex].Row == r &&
					searchResults[currentSearchIndex].Col == c {
					// Current match: vibrant cyan highlight
					backgroundColor = tcell.NewRGBColor(0, 180, 216)
					color = tcell.ColorBlack
					attributes = tcell.AttrBold
				} else {
					// Other matches: soft purple highlight
					backgroundColor = tcell.NewRGBColor(100, 100, 150)
					color = tcell.ColorWhite
					attributes = tcell.AttrNone
				}
			}

			// Determine max width for this column
			maxWidth := 0
			if width, isWrapped := wrappedColumns[c]; isWrapped {
				maxWidth = width
				// Truncate text if it exceeds max width
				cellText = truncateText(cellText, maxWidth)
			}

			// Create cell with modern styling
			cell := tview.NewTableCell(cellText).
				SetTextColor(color).
				SetBackgroundColor(backgroundColor).
				SetAttributes(attributes).
				SetAlign(alignment).
				SetExpansion(1)

			if maxWidth > 0 {
				cell.SetMaxWidth(maxWidth)
			}

			t.SetCell(r, c, cell)
		}
	}
}

// add stats data to stats table
func drawStats(s statsSummary, t *tview.Table) {
	t.Clear()
	summaryData := s.getSummaryData()
	rows, cols := len(summaryData), len(summaryData[0])

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			backgroundColor := tcell.ColorDefault

			// Modern styling: alternate row colors for better readability
			if r%2 == 1 {
				backgroundColor = tcell.NewRGBColor(20, 20, 30)
			}

			// Highlight stat labels with accent color
			if c == 0 {
				color = tcell.NewRGBColor(100, 200, 255) // Soft blue for labels
			}

			t.SetCell(r, c,
				tview.NewTableCell(summaryData[r][c]).
					SetTextColor(color).
					SetBackgroundColor(backgroundColor).
					SetAlign(tview.AlignLeft))
		}
	}
}

// draw app UI
func drawUI(b *Buffer) error {

	//bufferTable init with modern styling
	bufferTable = tview.NewTable()
	bufferTable.SetSelectable(true, true)
	bufferTable.SetBorders(false)
	bufferTable.SetSeparator(tview.Borders.Vertical)             // Add subtle vertical separators
	bufferTable.SetBordersColor(tcell.NewRGBColor(60, 100, 140)) // Subtle blue borders
	bufferTable.SetFixed(b.rowFreeze, b.colFreeze)
	bufferTable.Select(0, 0)
	bufferTable.SetSelectedStyle(tcell.Style{}.
		Foreground(tcell.ColorWhite).
		Background(tcell.NewRGBColor(80, 120, 160)). // Darker, muted blue
		Attributes(tcell.AttrBold))

	// Auto-detect and wrap long columns (sample first 100 rows, threshold 50 characters)
	detectAndWrapLongColumns(b, 100, 50)

	drawBuffer(b, bufferTable)

	//main page init with modern styling
	cursorPosStr = buildCursorPosStr(0, 0) //footer right
	if statusMessage == "" {
		statusMessage = "All Done"
	}
	shorFileName := filepath.Base(args.FileName)
	fileNameStr = shorFileName + "  |  " + "? help" //footer left
	filterInfoStr := buildFilterInfoStr(0)          // Top strip for filter info, initially at column 0

	mainPage = tview.NewFrame(bufferTable).
		SetBorders(0, 0, 0, 0, 0, 0)

	// Add filter info strip at top if filter is active and cursor on filtered column
	if filterInfoStr != "" {
		mainPage.AddText(filterInfoStr, true, tview.AlignCenter, tcell.NewRGBColor(255, 140, 0))
	}

	// Add main footer at bottom
	mainPage.AddText(fileNameStr, false, tview.AlignLeft, tcell.NewRGBColor(255, 150, 50)).
		AddText(statusMessage, false, tview.AlignCenter, tcell.NewRGBColor(100, 200, 255)).
		AddText(cursorPosStr, false, tview.AlignRight, tcell.NewRGBColor(150, 255, 150))

	drawFooterText := func(lstr, cstr, rstr string) {
		statusMessage = cstr // Update global status
		mainPage.Clear()

		// Add filter info strip at top if filter is active and cursor on filtered column
		filterInfoStr := buildFilterInfoStr(currentCursorColumn)
		if filterInfoStr != "" {
			mainPage.AddText(filterInfoStr, true, tview.AlignCenter, tcell.NewRGBColor(255, 140, 0))
		}

		// Add main footer at bottom
		mainPage.AddText(lstr, false, tview.AlignLeft, tcell.NewRGBColor(255, 150, 50)).
			AddText(cstr, false, tview.AlignCenter, tcell.NewRGBColor(100, 200, 255)).
			AddText(rstr, false, tview.AlignRight, tcell.NewRGBColor(150, 255, 150))
	}

	//UI init - add pages to UI container
	UI = tview.NewPages()
	UI.AddPage("main", mainPage, true, true)

	//bufferTable Event
	//bufferTable update cursor position
	bufferTable.SetSelectionChangedFunc(func(row int, column int) {
		// Mark that user has moved cursor if they moved from initial position (0,0)
		if !userMovedCursor && (row != 0 || column != 0) {
			userMovedCursor = true
		}

		// Update current cursor column
		currentCursorColumn = column

		cursorPosStr = buildCursorPosStr(row, column)

		// Rebuild the page with filter strip based on current column
		drawFooterText(fileNameStr, statusMessage, cursorPosStr)
	})

	//bufferTable HotKey Event
	bufferTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Mark that user is interacting with cursor movement keys
		if event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown ||
			event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyRight ||
			event.Key() == tcell.KeyHome || event.Key() == tcell.KeyEnd ||
			event.Key() == tcell.KeyPgUp || event.Key() == tcell.KeyPgDn ||
			(event.Key() == tcell.KeyRune && (event.Rune() == 'h' || event.Rune() == 'j' ||
				event.Rune() == 'k' || event.Rune() == 'l')) {
			userMovedCursor = true
		}

		// Vim-like navigation
		// h - move left
		if event.Key() == tcell.KeyRune && event.Rune() == 'h' {
			row, col := bufferTable.GetSelection()
			if col > 0 {
				bufferTable.Select(row, col-1)
			}
			return nil
		}

		// l - move right
		if event.Key() == tcell.KeyRune && event.Rune() == 'l' {
			row, col := bufferTable.GetSelection()
			if col < b.colLen-1 {
				bufferTable.Select(row, col+1)
			}
			return nil
		}

		// j - move down
		if event.Key() == tcell.KeyRune && event.Rune() == 'j' {
			row, col := bufferTable.GetSelection()
			if row < b.rowLen-1 {
				bufferTable.Select(row+1, col)
			}
			return nil
		}

		// k - move up
		if event.Key() == tcell.KeyRune && event.Rune() == 'k' {
			row, col := bufferTable.GetSelection()
			if row > 0 {
				bufferTable.Select(row-1, col)
			}
			return nil
		}

		// gg - go to first row
		if event.Key() == tcell.KeyRune && event.Rune() == 'g' {
			if lastKeyWasG {
				bufferTable.Select(0, 0)
				bufferTable.ScrollToBeginning()
				lastKeyWasG = false
				return nil
			}
			lastKeyWasG = true
			// Set a timer to reset lastKeyWasG after a short delay
			go func() {
				time.Sleep(500 * time.Millisecond)
				lastKeyWasG = false
			}()
			return nil
		}

		// G - go to last row
		if event.Key() == tcell.KeyRune && event.Rune() == 'G' {
			_, col := bufferTable.GetSelection()
			bufferTable.Select(b.rowLen-1, col)
			bufferTable.ScrollToEnd()
			return nil
		}

		// Ctrl+d - page down (half page)
		if event.Key() == tcell.KeyCtrlD {
			row, col := bufferTable.GetSelection()
			newRow := row + 10 // Move 10 rows down
			if newRow >= b.rowLen {
				newRow = b.rowLen - 1
			}
			bufferTable.Select(newRow, col)
			return nil
		}

		// Ctrl+u - page up (half page)
		if event.Key() == tcell.KeyCtrlU {
			row, col := bufferTable.GetSelection()
			newRow := row - 10 // Move 10 rows up
			if newRow < 0 {
				newRow = 0
			}
			bufferTable.Select(newRow, col)
			return nil
		}

		// 0 - go to first column
		if event.Key() == tcell.KeyRune && event.Rune() == '0' {
			row, _ := bufferTable.GetSelection()
			bufferTable.Select(row, 0)
			return nil
		}

		// $ - go to last column
		if event.Key() == tcell.KeyRune && event.Rune() == '$' {
			row, _ := bufferTable.GetSelection()
			bufferTable.Select(row, b.colLen-1)
			return nil
		}

		// w - move to next column (word forward)
		if event.Key() == tcell.KeyRune && event.Rune() == 'w' {
			row, col := bufferTable.GetSelection()
			if col < b.colLen-1 {
				bufferTable.Select(row, col+1)
			}
			return nil
		}

		// b - move to previous column (word backward)
		if event.Key() == tcell.KeyRune && event.Rune() == 'b' {
			row, col := bufferTable.GetSelection()
			if col > 0 {
				bufferTable.Select(row, col-1)
			}
			return nil
		}

		// / - search functionality
		if event.Key() == tcell.KeyRune && event.Rune() == '/' {
			// Create search form
			form := tview.NewForm()
			form.AddInputField("Search:", "", 40, nil, nil)
			form.AddCheckbox("Use Regex:", searchUseRegex, func(checked bool) {
				searchUseRegex = checked
			})
			form.AddCheckbox("Case Sensitive:", false, nil)
			form.GetFormItem(1).(*tview.Checkbox).SetLabelColor(tcell.NewRGBColor(180, 220, 220))
			form.GetFormItem(2).(*tview.Checkbox).SetLabelColor(tcell.NewRGBColor(180, 220, 220))
			form.GetFormItem(1).(*tview.Checkbox).SetFieldBackgroundColor(tcell.NewRGBColor(80, 80, 100)).SetFieldTextColor(tcell.NewRGBColor(0, 255, 255))
			form.GetFormItem(2).(*tview.Checkbox).SetFieldBackgroundColor(tcell.NewRGBColor(80, 80, 100)).SetFieldTextColor(tcell.NewRGBColor(0, 255, 255))

			// Define search execution function to avoid duplication
			executeSearch := func() {
				query := form.GetFormItem(0).(*tview.InputField).GetText()
				useRegex := form.GetFormItem(1).(*tview.Checkbox).IsChecked()
				caseSensitive := form.GetFormItem(2).(*tview.Checkbox).IsChecked()
				if query != "" {
					searchQuery = query
					searchUseRegex = useRegex
					searchResults = performSearch(b, query, useRegex, caseSensitive)

					if len(searchResults) > 0 {
						currentSearchIndex = 0
						bufferTable.Select(searchResults[0].Row, searchResults[0].Col)
						drawBuffer(b, bufferTable)
						searchMode := "matches"
						if useRegex {
							searchMode = "regex matches"
						}
						drawFooterText(fileNameStr,
							fmt.Sprintf("Found %d %s (1/%d)", len(searchResults), searchMode, len(searchResults)),
							cursorPosStr)
					} else {
						currentSearchIndex = -1
						if useRegex {
							drawFooterText(fileNameStr, "Invalid regex or no matches found", cursorPosStr)
						} else {
							drawFooterText(fileNameStr, "No matches found", cursorPosStr)
						}
					}
				}
				UI.HidePage("searchModal")
				app.SetFocus(bufferTable)
			}
			form.AddButton("Search", executeSearch)
			form.AddButton("Cancel", func() {
				UI.HidePage("searchModal")
				app.SetFocus(bufferTable)
			})
			form.SetButtonsAlign(tview.AlignCenter)
			form.SetBorder(true)
			title := " 🔍 Search - Tab to navigate, Enter to search, Esc to cancel "
			form.SetTitle(title)
			form.SetTitleAlign(tview.AlignCenter)
			form.SetBorderColor(tcell.NewRGBColor(0, 200, 255)) // Bright Blue
			form.SetBackgroundColor(tcell.NewRGBColor(20, 30, 40))
			form.SetLabelColor(tcell.NewRGBColor(180, 220, 220))
			form.SetFieldBackgroundColor(tcell.NewRGBColor(30, 40, 50))
			form.SetFieldTextColor(tcell.ColorWhite)
			form.SetButtonBackgroundColor(tcell.NewRGBColor(0, 200, 255))
			form.SetButtonTextColor(tcell.ColorBlack)
			searchButton := form.GetButton(0)
			searchButton.SetActivatedStyle(tcell.Style{}.
				Background(tcell.NewRGBColor(80, 120, 160)).
				Foreground(tcell.ColorWhite))
			cancelButton := form.GetButton(1)
			cancelButton.SetActivatedStyle(tcell.Style{}.
				Background(tcell.NewRGBColor(80, 120, 160)).
				Foreground(tcell.ColorWhite))

			// Handle Escape and Enter keys on form
			form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyEscape {
					UI.HidePage("searchModal")
					app.SetFocus(bufferTable)
					return nil
				}
				if event.Key() == tcell.KeyEnter {
					if itemIndex, _ := form.GetFocusedItemIndex(); itemIndex >= 0 {
						item := form.GetFormItem(itemIndex)
						if checkbox, ok := item.(*tview.Checkbox); ok {
							checkbox.SetChecked(!checkbox.IsChecked())
							return nil
						}
					}
					executeSearch()
					return nil
				}
				return event
			})

			// Create centered modal overlay
			searchModal = tview.NewFlex().
				AddItem(nil, 0, 1, false).
				AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(nil, 0, 1, false).
					AddItem(form, 11, 1, true).
					AddItem(nil, 0, 1, false), 60, 1, true).
				AddItem(nil, 0, 1, false)

			UI.AddPage("searchModal", searchModal, true, true)
			UI.ShowPage("searchModal")
			app.SetFocus(form)
			return nil
		}

		// Navigate to next search result
		if event.Key() == tcell.KeyRune && event.Rune() == 'n' {
			if len(searchResults) > 0 && currentSearchIndex >= 0 {
				currentSearchIndex = (currentSearchIndex + 1) % len(searchResults)
				bufferTable.Select(searchResults[currentSearchIndex].Row, searchResults[currentSearchIndex].Col)
				drawBuffer(b, bufferTable) // Redraw to update highlighting
				drawFooterText(fileNameStr,
					fmt.Sprintf("Match %d/%d", currentSearchIndex+1, len(searchResults)),
					cursorPosStr)
			} else if searchQuery != "" {
				drawFooterText(fileNameStr, "No search results. Press / to search", cursorPosStr)
			}
			return nil
		}

		// Navigate to previous search result
		if event.Key() == tcell.KeyRune && event.Rune() == 'N' {
			if len(searchResults) > 0 && currentSearchIndex >= 0 {
				currentSearchIndex--
				if currentSearchIndex < 0 {
					currentSearchIndex = len(searchResults) - 1
				}
				bufferTable.Select(searchResults[currentSearchIndex].Row, searchResults[currentSearchIndex].Col)
				drawBuffer(b, bufferTable) // Redraw to update highlighting
				drawFooterText(fileNameStr,
					fmt.Sprintf("Match %d/%d", currentSearchIndex+1, len(searchResults)),
					cursorPosStr)
			} else if searchQuery != "" {
				drawFooterText(fileNameStr, "No search results. Press / to search", cursorPosStr)
			}
			return nil
		}

		// Escape - clear search highlighting
		if event.Key() == tcell.KeyEscape {
			if searchQuery != "" {
				searchQuery = ""
				searchResults = []SearchResult{}
				currentSearchIndex = -1
				drawBuffer(b, bufferTable)
				drawFooterText(fileNameStr, "Search cleared", cursorPosStr)
			}
			return nil
		}

		// f - column filter functionality
		if event.Key() == tcell.KeyRune && event.Rune() == 'f' {
			_, column := bufferTable.GetSelection()

			// Create filter form
			filterForm := tview.NewForm()

			// Operator selection
			operators := []string{"contains", "equals", "starts with", "ends with", "regex", ">", "<", ">=", "<="}
			selectedOperatorIndex := 0

			// Value input
			query := ""
			caseSensitive := false

			if opts, exists := activeFilters[column]; exists {
				query = opts.Query
				caseSensitive = opts.CaseSensitive
				for i, op := range operators {
					if op == opts.Operator {
						selectedOperatorIndex = i
						break
					}
				}
			}

			filterForm.AddDropDown("Operator:", operators, selectedOperatorIndex, func(option string, optionIndex int) {
				selectedOperatorIndex = optionIndex
			})
			filterForm.AddInputField("Value:", query, 40, nil, nil)
			filterForm.AddCheckbox("Case Sensitive:", caseSensitive, func(checked bool) {
				caseSensitive = checked
			})
			filterForm.GetFormItem(2).(*tview.Checkbox).SetLabelColor(tcell.NewRGBColor(180, 220, 220))
			filterForm.GetFormItem(2).(*tview.Checkbox).SetFieldBackgroundColor(tcell.NewRGBColor(80, 80, 100)).SetFieldTextColor(tcell.NewRGBColor(0, 255, 255))

			applyFilter := func() {
				query = filterForm.GetFormItem(1).(*tview.InputField).GetText()
				operator := operators[selectedOperatorIndex]

				if query != "" {
					drawFooterText(fileNameStr, "Filtering...", cursorPosStr)
					app.ForceDraw()

					// Add or update filter for this column
					activeFilters[column] = FilterOptions{
						Query:         query,
						Operator:      operator,
						CaseSensitive: caseSensitive,
					}

					// Apply all filters starting from original buffer
					if originalBuffer == nil {
						originalBuffer = b // Save original buffer first time
					}

					// Start with original buffer and apply all filters sequentially
					filteredBuffer := originalBuffer
					for col, opts := range activeFilters {
						filteredBuffer = filteredBuffer.filterByColumn(col, opts)
					}

					// Update display with filtered data
					if filteredBuffer.rowLen <= filteredBuffer.rowFreeze {
						drawFooterText(fileNameStr, "No rows match filters", cursorPosStr)
						// Remove this filter since it results in no data
						delete(activeFilters, column)
					} else {
						// Replace current buffer with filtered buffer
						b = filteredBuffer
						isFiltered = true

						drawBuffer(b, bufferTable)
						bufferTable.Select(0, column) // Stay at same column, go to first row
						matchCount := b.rowLen - b.rowFreeze
						drawFooterText(fileNameStr,
							fmt.Sprintf("Filtered: %d rows match (%d filters active, r to reset)", matchCount, len(activeFilters)),
							cursorPosStr)
					}
				} else {
					// Empty query means remove filter for this column
					if _, exists := activeFilters[column]; exists {
						delete(activeFilters, column)

						// Reapply remaining filters
						if len(activeFilters) == 0 {
							// No more filters, restore original
							b = originalBuffer
							isFiltered = false
							drawBuffer(b, bufferTable)
							bufferTable.Select(0, column) // Stay at same column
							drawFooterText(fileNameStr, "All filters cleared - showing all rows", cursorPosStr)
						} else {
							// Apply remaining filters
							filteredBuffer := originalBuffer
							for col, opts := range activeFilters {
								filteredBuffer = filteredBuffer.filterByColumn(col, opts)
							}
							b = filteredBuffer
							drawBuffer(b, bufferTable)
							bufferTable.Select(0, column) // Stay at same column
							matchCount := b.rowLen - b.rowFreeze
							drawFooterText(fileNameStr,
								fmt.Sprintf("Filter removed: %d rows match (%d filters active)", matchCount, len(activeFilters)),
								cursorPosStr)
						}
					}
				}
				UI.HidePage("filterModal")
				app.SetFocus(bufferTable)
			}

			filterForm.AddButton("Filter", applyFilter)
			filterForm.AddButton("Cancel", func() {
				UI.HidePage("filterModal")
				app.SetFocus(bufferTable)
			})
			filterForm.SetButtonsAlign(tview.AlignCenter)
			filterForm.SetBorder(true)

			filterTitle := fmt.Sprintf(" 🔎 Filter Column %d - Enter to filter, Esc to cancel ", column)
			if _, exists := activeFilters[column]; exists {
				filterTitle = fmt.Sprintf(" 🔎 Edit Filter for Column %d (empty value to remove) - Enter to apply, Esc to cancel ", column)
			}
			filterForm.SetTitle(filterTitle)
			filterForm.SetTitleAlign(tview.AlignCenter)
			filterForm.SetBorderColor(tcell.NewRGBColor(0, 200, 255)) // Bright Blue
			filterForm.SetBackgroundColor(tcell.NewRGBColor(20, 30, 40))
			filterForm.SetLabelColor(tcell.NewRGBColor(180, 220, 220))
			filterForm.SetFieldBackgroundColor(tcell.NewRGBColor(30, 40, 50))
			filterForm.SetFieldTextColor(tcell.ColorWhite)
			filterForm.SetButtonBackgroundColor(tcell.NewRGBColor(0, 200, 255))
			filterForm.SetButtonTextColor(tcell.ColorBlack)
			filterButton := filterForm.GetButton(0)
			filterButton.SetActivatedStyle(tcell.Style{}.
				Background(tcell.NewRGBColor(80, 120, 160)).
				Foreground(tcell.ColorWhite))
			cancelButton := filterForm.GetButton(1)
			cancelButton.SetActivatedStyle(tcell.Style{}.
				Background(tcell.NewRGBColor(80, 120, 160)).
				Foreground(tcell.ColorWhite))

			// Handle Escape and Enter keys on form
			filterForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyEscape {
					UI.HidePage("filterModal")
					app.SetFocus(bufferTable)
					return nil
				}
				if event.Key() == tcell.KeyEnter {
					// If the dropdown has focus, let it handle the Enter key.
					if itemIndex, _ := filterForm.GetFocusedItemIndex(); itemIndex >= 0 {
						if item := filterForm.GetFormItem(itemIndex); item != nil {
							if _, ok := item.(*tview.DropDown); ok {
								return event
							}
							if checkbox, ok := item.(*tview.Checkbox); ok {
								checkbox.SetChecked(!checkbox.IsChecked())
								return nil
							}
						}
					}
					// if dropdown is open, pass enter to it
					if _, ok := app.GetFocus().(*tview.List); ok {
						return event
					}
					applyFilter()
					return nil
				}
				return event
			})

			// Create centered modal overlay
			filterModal := tview.NewFlex().
				AddItem(nil, 0, 1, false).
				AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(nil, 0, 1, false).
					AddItem(filterForm, 13, 1, true).
					AddItem(nil, 0, 1, false), 80, 1, true).
				AddItem(nil, 0, 1, false)

			UI.AddPage("filterModal", filterModal, true, true)
			UI.ShowPage("filterModal")
			app.SetFocus(filterForm)
			return nil
		}

		// r - reset filter for current column
		if event.Key() == tcell.KeyRune && event.Rune() == 'r' {
			if isFiltered && originalBuffer != nil {
				row, column := bufferTable.GetSelection()

				// Check if current column has a filter
				if _, hasFilter := activeFilters[column]; hasFilter {
					// Remove filter for this column
					delete(activeFilters, column)

					// Reapply remaining filters
					if len(activeFilters) == 0 {
						// No more filters, restore original
						b = originalBuffer
						isFiltered = false
						drawBuffer(b, bufferTable)
						bufferTable.Select(row, column)
						drawFooterText(fileNameStr, "All filters cleared - showing all rows", cursorPosStr)
					} else {
						// Apply remaining filters
						filteredBuffer := originalBuffer
						for col, opts := range activeFilters {
							filteredBuffer = filteredBuffer.filterByColumn(col, opts)
						}
						b = filteredBuffer
						drawBuffer(b, bufferTable)
						bufferTable.Select(row, column)
						matchCount := b.rowLen - b.rowFreeze
						drawFooterText(fileNameStr,
							fmt.Sprintf("Filter removed from current column: %d rows match (%d filters active)", matchCount, len(activeFilters)),
							cursorPosStr)
					}
				} else if len(activeFilters) > 0 {
					// Current column doesn't have a filter, but others do
					drawFooterText(fileNameStr, "Current column has no filter - navigate to filtered column to remove", cursorPosStr)
				}
			}
			return nil
		}

		// s - sort by column, ascending (s for sort)
		if event.Key() == tcell.KeyRune && event.Rune() == 's' {
			_, column := bufferTable.GetSelection()
			drawFooterText(fileNameStr, "Sorting...", cursorPosStr)
			app.ForceDraw()
			colType := b.getColType(column)
			switch colType {
			case colTypeFloat:
				b.sortByNum(column, false)
			case colTypeDate:
				b.sortByDate(column, false)
			default:
				b.sortByStr(column, false)
			}
			drawBuffer(b, bufferTable)
			drawFooterText(fileNameStr, "All Done", cursorPosStr)
		}

		// S - sort by column, descending (capital S for reverse sort)
		if event.Key() == tcell.KeyRune && event.Rune() == 'S' {
			_, column := bufferTable.GetSelection()
			drawFooterText(fileNameStr, "Sorting...", cursorPosStr)
			app.ForceDraw()
			colType := b.getColType(column)
			switch colType {
			case colTypeFloat:
				b.sortByNum(column, true)
			case colTypeDate:
				b.sortByDate(column, true)
			default:
				b.sortByStr(column, true)
			}
			drawBuffer(b, bufferTable)
			drawFooterText(fileNameStr, "All Done", cursorPosStr)
		}

		// i - show stats info for current column
		if event.Key() == tcell.KeyRune && event.Rune() == 'i' {
			_, column := bufferTable.GetSelection()
			drawFooterText(fileNameStr, "Calculating statistics...", cursorPosStr)
			app.ForceDraw()

			// Use the current buffer (which is filtered if filters are active)
			// This ensures stats are calculated only on visible/filtered data
			currentBuffer := b

			var statsS statsSummary
			summaryArray := currentBuffer.getCol(column)
			columnName := "Column " + I2S(column)

			// Get column name from header if available
			if currentBuffer.rowFreeze > 0 && len(currentBuffer.cont) > 0 && column < len(currentBuffer.cont[0]) {
				columnName = currentBuffer.cont[0][column]
				summaryArray = summaryArray[1:]
			}

			// Determine statistics type
			if currentBuffer.getColType(column) == colTypeFloat {
				statsS = &ContinuousStats{}
			} else {
				statsS = &DiscreteStats{}
			}
			statsS.summary(summaryArray)

			// Show statistics as a modal dialog with filter indication
			showStatsDialog(statsS, columnName, currentBuffer.getColType(column))
			drawFooterText(fileNameStr, "All Done", cursorPosStr)
			return nil
		}

		// t - toggle/change column data type (t for type)
		if event.Key() == tcell.KeyRune && event.Rune() == 't' {
			row, column := bufferTable.GetSelection()
			currentType := b.getColType(column)

			// Cycle through types: Str -> Num -> Date -> Str
			var newType int
			switch currentType {
			case colTypeStr:
				newType = colTypeFloat
			case colTypeFloat:
				newType = colTypeDate
			case colTypeDate:
				newType = colTypeStr
			default:
				newType = colTypeStr
			}

			b.setColType(column, newType)
			cursorPosStr = buildCursorPosStr(row, column)
			drawFooterText(fileNameStr, statusMessage, cursorPosStr)
		}

		// W - toggle text wrapping for current column (capital W for wrap)
		if event.Key() == tcell.KeyRune && event.Rune() == 'W' {
			_, column := bufferTable.GetSelection()

			if _, isWrapped := wrappedColumns[column]; isWrapped {
				// Unwrap: remove from wrapped columns
				delete(wrappedColumns, column)
				drawFooterText(fileNameStr, "Column width limit removed", cursorPosStr)
			} else {
				// Wrap: add to wrapped columns with default width
				width := getColumnMaxWidth(column)
				wrappedColumns[column] = width
				drawFooterText(fileNameStr, fmt.Sprintf("Column width limited to %d chars", width), cursorPosStr)
			}

			// Redraw the table with updated wrapping
			drawBuffer(b, bufferTable)
			return nil
		}

		// q - quit application
		if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			app.Stop()
			return nil
		}

		// ? - switch to help page
		if event.Key() == tcell.KeyRune && event.Rune() == '?' {
			showHelpDialog()
			return nil
		}

		app.ForceDraw()
		return event
	})
	// Add mouse handler for scrolling and clicking
	bufferTable.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		// Mark that user has interacted via mouse
		if action == tview.MouseLeftClick {
			userMovedCursor = true
		}

		// Handle mouse wheel scrolling
		switch action {
		case tview.MouseScrollUp:
			row, col := bufferTable.GetSelection()
			if row > 0 {
				bufferTable.Select(row-1, col)
			}
			return action, event
		case tview.MouseScrollDown:
			row, col := bufferTable.GetSelection()
			if row < b.rowLen-1 {
				bufferTable.Select(row+1, col)
			}
			return action, event
		}

		// Pass through other mouse events to default handler
		return action, event
	})

	//bufferTable quit event
	bufferTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
	})

	return nil
}

// showHelpDialog displays the help content as a centered modal dialog
func showHelpDialog() {
	// Create help content text view
	helpText := tview.NewTextView().
		SetDynamicColors(true).
		SetText(getHelpContent()).
		SetTextAlign(tview.AlignLeft).
		SetWordWrap(true)

	// Make help text scrollable with modern styling
	helpText.SetScrollable(true)
	helpText.SetBorder(true)
	helpText.SetTitle(" ❓ Help - Press ? or q or Esc to close ")
	helpText.SetTitleAlign(tview.AlignCenter)
	helpText.SetBorderColor(tcell.NewRGBColor(150, 100, 255))
	helpText.SetBackgroundColor(tcell.NewRGBColor(10, 10, 20))

	// Handle key events to close dialog
	helpText.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape ||
			(event.Key() == tcell.KeyRune && (event.Rune() == '?' || event.Rune() == 'q')) {
			UI.RemovePage("helpDialog")
			app.SetFocus(bufferTable)
			return nil
		}
		// Allow j/k navigation in help
		if event.Key() == tcell.KeyRune && event.Rune() == 'j' {
			row, col := helpText.GetScrollOffset()
			helpText.ScrollTo(row+1, col)
			return nil
		}
		if event.Key() == tcell.KeyRune && event.Rune() == 'k' {
			row, col := helpText.GetScrollOffset()
			if row > 0 {
				helpText.ScrollTo(row-1, col)
			}
			return nil
		}
		// gg - go to top
		if event.Key() == tcell.KeyRune && event.Rune() == 'g' {
			if lastKeyWasG {
				helpText.ScrollToBeginning()
				lastKeyWasG = false
				return nil
			}
			lastKeyWasG = true
			go func() {
				time.Sleep(500 * time.Millisecond)
				lastKeyWasG = false
			}()
			return nil
		}
		// G - go to bottom
		if event.Key() == tcell.KeyRune && event.Rune() == 'G' {
			helpText.ScrollToEnd()
			return nil
		}
		// Ctrl-d/u for page scrolling
		if event.Key() == tcell.KeyCtrlD {
			row, col := helpText.GetScrollOffset()
			helpText.ScrollTo(row+10, col)
			return nil
		}
		if event.Key() == tcell.KeyCtrlU {
			row, col := helpText.GetScrollOffset()
			if row > 10 {
				helpText.ScrollTo(row-10, col)
			} else {
				helpText.ScrollTo(0, col)
			}
			return nil
		}
		return event
	})

	// Create a centered modal with the help text
	// Modal dimensions: 80% width, 85% height
	helpModal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(helpText, 0, 85, true).
			AddItem(nil, 0, 1, false), 0, 80, true).
		AddItem(nil, 0, 1, false)

	// Add and show the help dialog
	UI.AddPage("helpDialog", helpModal, true, true)
	app.SetFocus(helpText)
}

// showStatsDialog displays column statistics as a centered modal dialog
func showStatsDialog(statsS statsSummary, columnName string, colType int) {
	// Create stats table
	statsTable := tview.NewTable()
	statsTable.SetSelectable(true, true)
	statsTable.SetBorders(false)
	statsTable.Select(0, 0)

	// Draw statistics
	drawStats(statsS, statsTable)

	// Create border with title and modern styling
	statsTable.SetBorder(true)
	statsTable.SetBorderColor(tcell.NewRGBColor(100, 200, 255))

	typeName := type2name(colType)
	title := fmt.Sprintf(" 📊 Statistics: %s [%s] ", columnName, typeName)

	// Add filter indicator if data is filtered
	if isFiltered && len(activeFilters) > 0 {
		title = fmt.Sprintf(" 📊 Statistics: %s [%s] (Filtered Data - %d filters active) ", columnName, typeName, len(activeFilters))
	}

	statsTable.SetTitle(title)
	statsTable.SetTitleAlign(tview.AlignCenter)

	// Create plot view
	plotView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(statsS.getPlot()).
		SetTextAlign(tview.AlignLeft)
	plotView.SetBorder(true)
	plotView.SetTitle(" 📈 Visual Distribution ")
	plotView.SetTitleAlign(tview.AlignCenter)
	plotView.SetBorderColor(tcell.NewRGBColor(255, 150, 50))

	// Create a flex layout with stats on left and plot on right
	statsContent := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(statsTable, 0, 1, true).
		AddItem(plotView, 0, 1, false)

	// Handle key events
	statsContent.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Escape or q - close dialog
		if event.Key() == tcell.KeyEscape || (event.Key() == tcell.KeyRune && event.Rune() == 'q') {
			UI.RemovePage("statsDialog")
			app.SetFocus(bufferTable)
			return nil
		}

		// gg - go to top
		if event.Key() == tcell.KeyRune && event.Rune() == 'g' {
			if lastKeyWasG {
				_, column := statsTable.GetSelection()
				statsTable.Select(0, column)
				statsTable.ScrollToBeginning()
				lastKeyWasG = false
				return nil
			}
			lastKeyWasG = true
			go func() {
				time.Sleep(500 * time.Millisecond)
				lastKeyWasG = false
			}()
			return nil
		}

		// G - go to bottom
		if event.Key() == tcell.KeyRune && event.Rune() == 'G' {
			_, column := statsTable.GetSelection()
			statsTable.Select(statsTable.GetRowCount()-1, column)
			statsTable.ScrollToEnd()
			return nil
		}

		// j/k navigation
		if event.Key() == tcell.KeyRune && event.Rune() == 'j' {
			row, col := statsTable.GetSelection()
			if row < statsTable.GetRowCount()-1 {
				statsTable.Select(row+1, col)
			}
			return nil
		}
		if event.Key() == tcell.KeyRune && event.Rune() == 'k' {
			row, col := statsTable.GetSelection()
			if row > 0 {
				statsTable.Select(row-1, col)
			}
			return nil
		}

		// Ctrl-d/u for page scrolling
		if event.Key() == tcell.KeyCtrlD {
			row, col := statsTable.GetSelection()
			newRow := row + 10
			if newRow >= statsTable.GetRowCount() {
				newRow = statsTable.GetRowCount() - 1
			}
			statsTable.Select(newRow, col)
			return nil
		}
		if event.Key() == tcell.KeyCtrlU {
			row, col := statsTable.GetSelection()
			newRow := row - 10
			if newRow < 0 {
				newRow = 0
			}
			statsTable.Select(newRow, col)
			return nil
		}

		return event
	})

	// Create a centered modal with the stats content
	// Modal dimensions: 80% width, 80% height
	statsModal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(statsContent, 0, 80, true).
			AddItem(tview.NewTextView().
				SetText("Press q or Esc to close").
				SetTextAlign(tview.AlignCenter).
				SetTextColor(tcell.NewRGBColor(150, 150, 150)), 1, 0, false).
			AddItem(nil, 0, 1, false), 0, 80, true).
		AddItem(nil, 0, 1, false)

	// Add and show the stats dialog
	UI.AddPage("statsDialog", statsModal, true, true)
	app.SetFocus(statsContent)
}
