package main

import (
	"path/filepath"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

//add buffer data to buffer table
func drawBuffer(b *Buffer, t *tview.Table, trs bool) {
	t.Clear()
	if trs {
		b.transpose()
	}
	cols, rows := b.colLen, b.rowLen

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			if c < b.colFreeze || r < b.rowFreeze {
				color = tcell.ColorYellow
			}
			if r == 0 && args.Header != -1 && args.Header != 2 {
				t.SetCell(r, c,
					tview.NewTableCell(b.cont[r][c]).
						SetTextColor(color).
						SetAlign(tview.AlignLeft))
				continue
			}
			t.SetCell(r, c,
				tview.NewTableCell(b.cont[r][c]).
					SetTextColor(color).
					SetAlign(tview.AlignLeft))
		}
	}
}

//add stats data to stats table
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

//draw app UI
func drawUI(b *Buffer, trs bool) error {

	//bufferTable init
	bufferTable := tview.NewTable()
	bufferTable.SetSelectable(true, true)
	bufferTable.SetBorders(false)
	bufferTable.SetFixed(b.rowFreeze, b.colFreeze)
	bufferTable.Select(0, 0)
	drawBuffer(b, bufferTable, trs)

	//main page init
	cursorPosStr := "Column Type: " + type2name(b.getColType(0)) + "  |  0,0  " //footer right
	infoStr := "All Done"                                                       //footer middle
	shorFileName := filepath.Base(args.FileName)
	fileNameStr := shorFileName + "  |  " + "? help page" //footer left
	mainPage := tview.NewFrame(bufferTable).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(fileNameStr, false, tview.AlignLeft, tcell.ColorDarkOrange).
		AddText(infoStr, false, tview.AlignCenter, tcell.ColorDarkOrange).
		AddText(cursorPosStr, false, tview.AlignRight, tcell.ColorDarkOrange)

	drawFooterText := func(lstr, cstr, rstr string) {
		mainPage.Clear()
		mainPage = mainPage.
			AddText(lstr, false, tview.AlignLeft, tcell.ColorDarkOrange).
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
	UI.AddPage("help", helpPage, true, true)
	UI.AddPage("stats", statsPage, true, true)
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
		cursorPosStr = "Column Type: " + type2name(b.getColType(column)) + "  |  " + strconv.Itoa(row) + "," + strconv.Itoa(column) + "  "
		drawFooterText(fileNameStr, infoStr, cursorPosStr)
	})

	//bufferTable HotKey Event
	bufferTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		//sort by column, ascend
		if event.Key() == tcell.KeyCtrlK {
			_, column := bufferTable.GetSelection()
			infoStr = "Sorting..."
			drawFooterText(fileNameStr, infoStr, cursorPosStr)
			app.ForceDraw()
			if b.getColType(column) == colTypeFloat {
				b.sortByNum(column, false)
			} else {
				b.sortByStr(column, false)
			}
			drawBuffer(b, bufferTable, trs)
			infoStr = "All Done"
			drawFooterText(fileNameStr, infoStr, cursorPosStr)
		}
		//sort by column, descend
		if event.Key() == tcell.KeyCtrlL {
			_, column := bufferTable.GetSelection()
			infoStr = "Sorting..."
			drawFooterText(fileNameStr, infoStr, cursorPosStr)
			app.ForceDraw()
			if b.getColType(column) == colTypeFloat {
				b.sortByNum(column, true)
			} else {
				b.sortByStr(column, true)
			}
			drawBuffer(b, bufferTable, trs)
			infoStr = "All Done"
			drawFooterText(fileNameStr, infoStr, cursorPosStr)
		}

		//show current column's stats info
		if event.Key() == tcell.KeyCtrlY {
			_, column := bufferTable.GetSelection()
			infoStr = "Calculating"
			drawFooterText(fileNameStr, infoStr, cursorPosStr)
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
			statsTable.Select(0, 0).GetFocusable()
			app.SetFocus(statsTable)
			statsTable.ScrollToBeginning()
			drawStats(statsS, statsTable)
			UI.SwitchToPage("stats")
			infoStr = "All Done"
			drawFooterText(fileNameStr, infoStr, cursorPosStr)
		}

		//change column data type
		if event.Key() == tcell.KeyCtrlM {
			row, column := bufferTable.GetSelection()
			var colType int
			if b.getColType(column) == colTypeFloat {
				colType = colTypeStr
			} else {
				colType = colTypeFloat
			}

			b.setColType(column, colType)
			cursorPosStr = "Column Type: " + type2name(b.getColType(column)) + "  |  " + strconv.Itoa(row) + "," + strconv.Itoa(column) + "  "
			drawFooterText(fileNameStr, infoStr, cursorPosStr)
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
