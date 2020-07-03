package main

import (
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
					tview.NewTableCell("("+type2name(b.colType[c])+")"+b.cont[r][c]).
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
	cursorPosStr := "0,0"
	mainPage := tview.NewFrame(bufferTable).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(args.FileName, false, tview.AlignLeft, tcell.ColorDarkOrange).
		AddText(cursorPosStr, false, tview.AlignRight, tcell.ColorDarkOrange)

	//statsTable init
	statsTable := tview.NewTable()
	statsTable.SetSelectable(true, true)
	statsTable.SetBorders(false)
	statsTable.Select(0, 0)

	// stats page init
	statsPage := tview.NewFrame(statsTable).
		SetBorders(0, 0, 0, 1, 0, 0).
		AddText("Basic Stats", true, tview.AlignCenter, tcell.ColorDarkOrange)

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
			statsTable.Select(0,column)
			statsTable.ScrollToBeginning()
		}

		//go to end of current column
		if event.Key() == tcell.KeyCtrlE {
			_, column := statsTable.GetSelection()
			statsTable.Select(statsTable.GetRowCount()-1,column)
			statsTable.ScrollToEnd()
		}
		return event
	})

	statsTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
	})

	//UI init
	UI = tview.NewPages()
	UI.AddPage("stats", statsPage, true, true)
	UI.AddPage("main", mainPage, true, true)

	//bufferTable Event
	//bufferTable update cursor postion
	bufferTable.SetSelectionChangedFunc(func(row int, column int) {
		cursorPosStr = strconv.Itoa(row) + "," + strconv.Itoa(column)
		mainPage.Clear()
		mainPage.AddText(args.FileName, false, tview.AlignLeft, tcell.ColorDarkOrange).
			AddText(cursorPosStr, false, tview.AlignRight, tcell.ColorDarkOrange)

	})

	//bufferTable HotKey Event
	bufferTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		//sort by column, ascend
		if event.Key() == tcell.KeyCtrlK {
			_, column := bufferTable.GetSelection()
			if b.getColType(column) == colTypeFloat {
				b.sortByNum(column, false)
			} else {
				b.sortByStr(column, false)
			}
			drawBuffer(b, bufferTable, trs)
		}
		//sort by column, descend
		if event.Key() == tcell.KeyCtrlL {
			_, column := bufferTable.GetSelection()
			if b.getColType(column) == colTypeFloat {
				b.sortByNum(column, true)
			} else {
				b.sortByStr(column, true)
			}
			drawBuffer(b, bufferTable, trs)
		}

		//show current column's stats info
		if event.Key() == tcell.KeyCtrlY {
			_, column := bufferTable.GetSelection()
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
			UI.SwitchToPage("stats")
			app.SetFocus(statsTable)
			statsTable.ScrollToBeginning()
			drawStats(statsS, statsTable)
		}

		//change column data type
		if event.Key() == tcell.KeyCtrlM {
			_, column := bufferTable.GetSelection()
			var colType int
			if b.getColType(column) == colTypeFloat{
				colType = colTypeStr
			}else {
				colType = colTypeFloat
			}

			b.setColType(column, colType)
			drawBuffer(b, bufferTable, trs)
		}

		//go to head of current column
		if event.Key() == tcell.KeyCtrlH {
			_, column := bufferTable.GetSelection()
			bufferTable.Select(0,column)
			bufferTable.ScrollToBeginning()
			drawBuffer(b, bufferTable, trs)
		}

		//go to end of current column
		if event.Key() == tcell.KeyCtrlE {
			_, column := bufferTable.GetSelection()
			bufferTable.Select(b.rowLen-1,column)
			bufferTable.ScrollToEnd()
			drawBuffer(b, bufferTable, trs)
		}
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
