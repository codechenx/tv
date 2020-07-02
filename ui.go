package main

import (
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func drawBuffer(b *Buffer, t *tview.Table, trs bool) {
	cols, rows := b.colLen, b.rowLen

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			if c < b.colFreeze || r < b.rowFreeze {
				color = tcell.ColorYellow
			}
			// check transpose view
			if !trs {
				t.SetCell(r, c,
					tview.NewTableCell(b.cont[r][c]).
						SetTextColor(color).
						SetAlign(tview.AlignLeft))
			} else {
				t.SetCell(c, r,
					tview.NewTableCell(b.cont[r][c]).
						SetTextColor(color).
						SetAlign(tview.AlignLeft))
			}
		}
	}
}

func drawUI(b *Buffer, trs bool) error {

	//table init
	table := tview.NewTable()
	drawBuffer(b, table, trs)
	table.SetSelectable(true, true)
	table.SetBorders(false)
	table.SetFixed(b.rowFreeze, b.colFreeze)

	//pages
	cursorPosStr := "0,0"
	mainPage := tview.NewFrame(table).
		SetBorders(0, 0, 0, 1, 0, 0).
		AddText(args.FileName, false, tview.AlignLeft, tcell.ColorDarkOrange).
		AddText(cursorPosStr, false, tview.AlignRight, tcell.ColorDarkOrange)

	statsPage := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText("OK")

	//UI init
	UI = tview.NewPages()
	UI.AddPage("stats", statsPage, true, true)
	UI.AddPage("main", mainPage, true, true)

	//mainPage Input Key Event
	//key envent - set cursor postion
	table.SetSelectionChangedFunc(func(row int, column int) {
		cursorPosStr = strconv.Itoa(row) + "," + strconv.Itoa(column)
		mainPage.Clear()
		mainPage.AddText(args.FileName, false, tview.AlignLeft, tcell.ColorDarkOrange).
			AddText(cursorPosStr, false, tview.AlignRight, tcell.ColorDarkOrange)

	})
	//key event -- sort data
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlS {
			_, column := table.GetSelection()
			b.sort(column, false)
			drawBuffer(b, table, trs)
		}

		if event.Key() == tcell.KeyCtrlL {
			_, column := table.GetSelection()
			b.sort(column, true)
			drawBuffer(b, table, trs)
		}

		if event.Key() == tcell.KeyCtrlX {
			UI.SwitchToPage("stats")
		}
		return event
	})

	//key event -- mark
	table.SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
	})

	//key event -- quit
	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
	})

	//statsPage Input Key Event
	statsPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlX {
			UI.SwitchToPage("main")
		}
		return event
	})

	return nil
}
