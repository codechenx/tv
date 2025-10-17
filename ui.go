package main

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// add buffer data to buffer table
func drawBuffer(b *Buffer, t *tview.Table, trs bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	t.Clear()
	if trs {
		b.transpose()
	}
	cols, rows := b.colLen, b.rowLen

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			backgroundColor := tcell.ColorDefault
			if c < b.colFreeze || r < b.rowFreeze {
				color = tcell.ColorYellow
			}

			// Get cell content
			cellText := b.cont[r][c]
			
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
			
			// Highlight search matches
			if isSearchMatch {
				// Check if this is the current search result
				if currentSearchIndex >= 0 && currentSearchIndex < len(searchResults) &&
					searchResults[currentSearchIndex].Row == r && 
					searchResults[currentSearchIndex].Col == c {
					// Current match: bright highlight
					backgroundColor = tcell.ColorDarkCyan
					color = tcell.ColorBlack
				} else {
					// Other matches: subtle highlight
					backgroundColor = tcell.ColorDarkGray
					color = tcell.ColorWhite
				}
			}

			// Apply text wrapping if column is marked for wrapping
			if maxWidth, isWrapped := wrappedColumns[c]; isWrapped {
				cellText = wrapText(cellText, maxWidth)
			}

			if r == 0 && args.Header != -1 && args.Header != 2 {
				t.SetCell(r, c,
					tview.NewTableCell(cellText).
						SetTextColor(color).
						SetBackgroundColor(backgroundColor).
						SetAlign(tview.AlignLeft).
						SetMaxWidth(0). // 0 means no limit, allows wrapping
						SetExpansion(1))
				continue
			}
			t.SetCell(r, c,
				tview.NewTableCell(cellText).
					SetTextColor(color).
					SetBackgroundColor(backgroundColor).
					SetAlign(tview.AlignLeft).
					SetMaxWidth(0).
					SetExpansion(1))
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
			t.SetCell(r, c,
				tview.NewTableCell(summaryData[r][c]).
					SetTextColor(color).
					SetAlign(tview.AlignLeft))
		}
	}
}

// draw app UI
func drawUI(b *Buffer, trs bool) error {

	//bufferTable init
	bufferTable = tview.NewTable()
	bufferTable.SetSelectable(true, true)
	bufferTable.SetBorders(false)
	bufferTable.SetFixed(b.rowFreeze, b.colFreeze)
	bufferTable.Select(0, 0)
	drawBuffer(b, bufferTable, trs)

	//main page init
	cursorPosStr = "Column Type: " + type2name(b.getColType(0)) + "  |  0,0  " //footer right
	if statusMessage == "" {
		statusMessage = "All Done"
	}
	shorFileName := filepath.Base(args.FileName)
	fileNameStr = shorFileName + "  |  " + "? help page" //footer left
	mainPage = tview.NewFrame(bufferTable).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(fileNameStr, false, tview.AlignLeft, tcell.ColorDarkOrange).
		AddText(statusMessage, false, tview.AlignCenter, tcell.ColorDarkOrange).
		AddText(cursorPosStr, false, tview.AlignRight, tcell.ColorDarkOrange)

	drawFooterText := func(lstr, cstr, rstr string) {
		statusMessage = cstr // Update global status
		mainPage.Clear()
		mainPage.AddText(lstr, false, tview.AlignLeft, tcell.ColorDarkOrange).
			AddText(cstr, false, tview.AlignCenter, tcell.ColorDarkOrange).
			AddText(rstr, false, tview.AlignRight, tcell.ColorDarkOrange)
	}
	//statsTable init
	statsTable := tview.NewTable()
	statsTable.SetSelectable(true, true)
	statsTable.SetBorders(false)
	statsTable.Select(0, 0)

	// stats page init
	statsPage := tview.NewFrame(statsTable).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Basic Stats", true, tview.AlignCenter, tcell.ColorDarkOrange)
	//help page init
	helpPage := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).SetText(getHelpContent())
	//UI init
	UI = tview.NewPages()
	UI.AddPage("help", helpPage, true, false)
	UI.AddPage("stats", statsPage, true, false)
	UI.AddPage("main", mainPage, true, true)

	//helpPage HotKey Event
	helpPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == '?' {
			UI.SwitchToPage("main")
		}
		return event
	})

	//statsPage HotKey Event
	statsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		//back to main page
		if event.Key() == tcell.KeyCtrlY {
			UI.SwitchToPage("main")
			app.SetFocus(bufferTable)
		}

		//go to head of current column
		if event.Key() == tcell.KeyCtrlH {
			_, column := statsTable.GetSelection()
			statsTable.Select(0, column)
			statsTable.ScrollToBeginning()
		}

		//go to end of current column
		if event.Key() == tcell.KeyCtrlE {
			_, column := statsTable.GetSelection()
			statsTable.Select(statsTable.GetRowCount()-1, column)
			statsTable.ScrollToEnd()
		}
		return event

	})

	statsTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
	})

	//bufferTable Event
	//bufferTable update cursor position
	bufferTable.SetSelectionChangedFunc(func(row int, column int) {
		// Mark that user has moved cursor if they moved from initial position (0,0)
		if !userMovedCursor && (row != 0 || column != 0) {
			userMovedCursor = true
		}
		cursorPosStr = "Column Type: " + type2name(b.getColType(column)) + "  |  " + strconv.Itoa(row) + "," + strconv.Itoa(column) + "  "
		drawFooterText(fileNameStr, statusMessage, cursorPosStr)
	})

	//bufferTable HotKey Event
	bufferTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Mark that user is interacting with cursor movement keys
		if event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown ||
			event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyRight ||
			event.Key() == tcell.KeyHome || event.Key() == tcell.KeyEnd ||
			event.Key() == tcell.KeyPgUp || event.Key() == tcell.KeyPgDn {
			userMovedCursor = true
		}

		//search functionality
		if event.Key() == tcell.KeyRune && event.Rune() == '/' {
			// Create search form
			var form *tview.Form
			form = tview.NewForm()
			form.AddInputField("Search:", "", 40, nil, nil)
			form.AddButton("Search", func() {
				query := form.GetFormItem(0).(*tview.InputField).GetText()
				if query != "" {
					searchQuery = query
					searchResults = performSearch(b, query, false)
					
					if len(searchResults) > 0 {
						currentSearchIndex = 0
						bufferTable.Select(searchResults[0].Row, searchResults[0].Col)
						drawBuffer(b, bufferTable, args.Transpose)
						drawFooterText(fileNameStr, 
							fmt.Sprintf("Found %d matches (1/%d)", len(searchResults), len(searchResults)), 
							cursorPosStr)
					} else {
						currentSearchIndex = -1
						drawFooterText(fileNameStr, "No matches found", cursorPosStr)
					}
				}
				UI.HidePage("searchModal")
				app.SetFocus(bufferTable)
			})
			form.AddButton("Cancel", func() {
				UI.HidePage("searchModal")
				app.SetFocus(bufferTable)
			})
			form.SetButtonsAlign(tview.AlignCenter)
			form.SetBorder(true)
			form.SetTitle(" Search (case-insensitive) - Enter to search, Esc to cancel ")
			form.SetTitleAlign(tview.AlignCenter)
			
			// Handle Escape and Enter keys on form
			form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyEscape {
					UI.HidePage("searchModal")
					app.SetFocus(bufferTable)
					return nil
				}
				if event.Key() == tcell.KeyEnter {
					query := form.GetFormItem(0).(*tview.InputField).GetText()
					if query != "" {
						searchQuery = query
						searchResults = performSearch(b, query, false)
						
						if len(searchResults) > 0 {
							currentSearchIndex = 0
							bufferTable.Select(searchResults[0].Row, searchResults[0].Col)
							drawBuffer(b, bufferTable, args.Transpose)
							drawFooterText(fileNameStr, 
								fmt.Sprintf("Found %d matches (1/%d)", len(searchResults), len(searchResults)), 
								cursorPosStr)
						} else {
							currentSearchIndex = -1
							drawFooterText(fileNameStr, "No matches found", cursorPosStr)
						}
					}
					UI.HidePage("searchModal")
					app.SetFocus(bufferTable)
					return nil
				}
				return event
			})
			
			// Create centered modal overlay
			searchModal = tview.NewFlex().
				AddItem(nil, 0, 1, false).
				AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(nil, 0, 1, false).
					AddItem(form, 9, 1, true).
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
				drawBuffer(b, bufferTable, args.Transpose) // Redraw to update highlighting
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
				drawBuffer(b, bufferTable, args.Transpose) // Redraw to update highlighting
				drawFooterText(fileNameStr, 
					fmt.Sprintf("Match %d/%d", currentSearchIndex+1, len(searchResults)), 
					cursorPosStr)
			} else if searchQuery != "" {
				drawFooterText(fileNameStr, "No search results. Press / to search", cursorPosStr)
			}
			return nil
		}

		// Clear search highlighting (Ctrl+/)
		if event.Key() == tcell.KeyCtrlUnderscore { // Ctrl+/ is often mapped to Ctrl+_
			if searchQuery != "" {
				searchQuery = ""
				searchResults = []SearchResult{}
				currentSearchIndex = -1
				drawBuffer(b, bufferTable, args.Transpose)
				drawFooterText(fileNameStr, "Search cleared", cursorPosStr)
			}
			return nil
		}

		// Column filter functionality (f key)
		if event.Key() == tcell.KeyRune && event.Rune() == 'f' {
			_, column := bufferTable.GetSelection()
			
			// Create filter form
			var filterForm *tview.Form
			filterForm = tview.NewForm()
			filterForm.AddInputField("Filter column by value:", "", 40, nil, nil)
			filterForm.AddButton("Filter", func() {
				query := filterForm.GetFormItem(0).(*tview.InputField).GetText()
				if query != "" {
					drawFooterText(fileNameStr, "Filtering...", cursorPosStr)
					app.ForceDraw()
					
					// Apply filter
					filteredBuffer := b.filterByColumn(column, query, false)
					
					// Update display with filtered data
					if filteredBuffer.rowLen <= filteredBuffer.rowFreeze {
						drawFooterText(fileNameStr, "No rows match filter", cursorPosStr)
					} else {
						// Replace current buffer with filtered buffer
						originalBuffer = b // Save original buffer
						b = filteredBuffer
						isFiltered = true
						filterColumn = column
						filterQuery = query
						
						drawBuffer(b, bufferTable, args.Transpose)
						bufferTable.Select(0, 0)
						matchCount := b.rowLen - b.rowFreeze
						drawFooterText(fileNameStr, 
							fmt.Sprintf("Filtered: %d rows match (Ctrl+R to reset)", matchCount), 
							cursorPosStr)
					}
				}
				UI.HidePage("filterModal")
				app.SetFocus(bufferTable)
			})
			filterForm.AddButton("Cancel", func() {
				UI.HidePage("filterModal")
				app.SetFocus(bufferTable)
			})
			filterForm.SetButtonsAlign(tview.AlignCenter)
			filterForm.SetBorder(true)
			filterForm.SetTitle(fmt.Sprintf(" Filter Column %d (case-insensitive) - Enter to filter, Esc to cancel ", column))
			filterForm.SetTitleAlign(tview.AlignCenter)
			
			// Handle Escape and Enter keys on form
			filterForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyEscape {
					UI.HidePage("filterModal")
					app.SetFocus(bufferTable)
					return nil
				}
				if event.Key() == tcell.KeyEnter {
					query := filterForm.GetFormItem(0).(*tview.InputField).GetText()
					if query != "" {
						drawFooterText(fileNameStr, "Filtering...", cursorPosStr)
						app.ForceDraw()
						
						// Apply filter
						filteredBuffer := b.filterByColumn(column, query, false)
						
						// Update display with filtered data
						if filteredBuffer.rowLen <= filteredBuffer.rowFreeze {
							drawFooterText(fileNameStr, "No rows match filter", cursorPosStr)
						} else {
							// Replace current buffer with filtered buffer
							originalBuffer = b // Save original buffer
							b = filteredBuffer
							isFiltered = true
							filterColumn = column
							filterQuery = query
							
							drawBuffer(b, bufferTable, args.Transpose)
							bufferTable.Select(0, 0)
							matchCount := b.rowLen - b.rowFreeze
							drawFooterText(fileNameStr, 
								fmt.Sprintf("Filtered: %d rows match (Ctrl+R to reset)", matchCount), 
								cursorPosStr)
						}
					}
					UI.HidePage("filterModal")
					app.SetFocus(bufferTable)
					return nil
				}
				return event
			})
			
			// Create centered modal overlay
			filterModal := tview.NewFlex().
				AddItem(nil, 0, 1, false).
				AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(nil, 0, 1, false).
					AddItem(filterForm, 9, 1, true).
					AddItem(nil, 0, 1, false), 70, 1, true).
				AddItem(nil, 0, 1, false)
			
			UI.AddPage("filterModal", filterModal, true, true)
			UI.ShowPage("filterModal")
			app.SetFocus(filterForm)
			return nil
		}

		// Reset filter (Ctrl+R)
		if event.Key() == tcell.KeyCtrlR {
			if isFiltered && originalBuffer != nil {
				b = originalBuffer
				isFiltered = false
				drawBuffer(b, bufferTable, args.Transpose)
				bufferTable.Select(0, 0)
				drawFooterText(fileNameStr, "Filter cleared - showing all rows", cursorPosStr)
			}
			return nil
		}

		//sort by column, ascend
		if event.Key() == tcell.KeyCtrlK {
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
			drawBuffer(b, bufferTable, trs)
			drawFooterText(fileNameStr, "All Done", cursorPosStr)
		}
		//sort by column, descend
		if event.Key() == tcell.KeyCtrlL {
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
			drawBuffer(b, bufferTable, trs)
			drawFooterText(fileNameStr, "All Done", cursorPosStr)
		}

		//show current column's stats info
		if event.Key() == tcell.KeyCtrlY {
			_, column := bufferTable.GetSelection()
			drawFooterText(fileNameStr, "Calculating", cursorPosStr)
			var statsS statsSummary
			summaryArray := b.getCol(column)
			if I2B(b.colFreeze) {
				summaryArray = summaryArray[1:]
			}
			if b.getColType(column) == colTypeFloat {
				statsS = &ContinuousStats{}
			} else {
				statsS = &DiscreteStats{}
			}
			statsS.summary(summaryArray)
			statsTable.Select(0, 0)
			app.SetFocus(statsTable)
			statsTable.ScrollToBeginning()
			drawStats(statsS, statsTable)
			UI.SwitchToPage("stats")
			drawFooterText(fileNameStr, "All Done", cursorPosStr)
		}

		//change column data type (cycle through Str -> Num -> Date -> Str)
		if event.Key() == tcell.KeyCtrlM {
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
			cursorPosStr = "Column Type: " + type2name(b.getColType(column)) + "  |  " + strconv.Itoa(row) + "," + strconv.Itoa(column) + "  "
			drawFooterText(fileNameStr, statusMessage, cursorPosStr)
		}

		//toggle text wrapping for current column
		if event.Key() == tcell.KeyCtrlW {
			_, column := bufferTable.GetSelection()

			if _, isWrapped := wrappedColumns[column]; isWrapped {
				// Unwrap: remove from wrapped columns
				delete(wrappedColumns, column)
			} else {
				// Wrap: add to wrapped columns with default width
				wrappedColumns[column] = 25 // Default wrap width (25 characters)
			}

			// Redraw the table with updated wrapping
			drawBuffer(b, bufferTable, args.Transpose)
		}

		//go to head of current column
		if event.Key() == tcell.KeyCtrlH {
			_, column := bufferTable.GetSelection()
			bufferTable.Select(0, column)
			bufferTable.ScrollToBeginning()
		}

		//go to end of current column
		if event.Key() == tcell.KeyCtrlE {
			_, column := bufferTable.GetSelection()
			bufferTable.Select(b.rowLen-1, column)
			bufferTable.ScrollToEnd()
		}
		//switch to help page
		if event.Key() == tcell.KeyRune && event.Rune() == '?' {
			UI.SwitchToPage("help")
		}

		if event.Key() == tcell.KeyRune && event.Rune() == 'G' {
			bufferTable.Select(b.rowLen-1, b.colLen-1)
		}

		if event.Key() == tcell.KeyRune && event.Rune() == 'g' {
			bufferTable.Select(0, 0)
			bufferTable.ScrollToBeginning()
		}
		app.ForceDraw()
		return event
	})

	//bufferTable quit event
	bufferTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
	})

	return nil
}
