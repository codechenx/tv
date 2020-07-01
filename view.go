package main

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func render(b *Buffer, trs bool) error {
	var fixPos [2]int
	switch args.Header {
	case -1:
		fixPos = [...]int{0, 0}
	case 0:
		fixPos = [...]int{1, 1}
	case 1:
		fixPos = [...]int{1, 0}
	case 2:
		fixPos = [...]int{0, 1}

	}

	cols, rows := b.colNum, b.rowNum
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			if c < fixPos[1] || r < fixPos[0] {
				color = tcell.ColorYellow
			}
			// check transpose view
			if !trs {
				table.SetCell(r, c,
					tview.NewTableCell(b.cont[r][c]).
						SetTextColor(color).
						SetAlign(tview.AlignLeft))
			} else {
				table.SetCell(c, r,
					tview.NewTableCell(b.cont[r][c]).
						SetTextColor(color).
						SetAlign(tview.AlignLeft))
			}
		}
	}

	footerPos := tview.NewTextView().
		SetTextAlign(tview.AlignRight).
		SetText("0,0")
	footerFP := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText(args.FileName)
	sepLine := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(strings.Repeat("-", 10000))
	footergrid := tview.NewGrid().SetColumns(0, 1).
		SetBorders(false).
		AddItem(footerFP, 0, 0, 1, 3, 0, 0, false).
		AddItem(footerPos, 0, 1, 1, 3, 0, 0, false)
	grid = tview.NewGrid().
		SetRows(0, 1, 1).
		SetBorders(false).
		AddItem(table, 0, 0, 1, 3, 0, 0, true).
		AddItem(sepLine, 1, 0, 1, 3, 0, 0, true).
		AddItem(footergrid, 2, 0, 1, 3, 0, 0, false)

	table.SetSelectable(true, true)
	table.SetBorders(false)
	//display cursor postion
	table.SetSelectionChangedFunc(func(row int, column int) {
		footerPos.SetText(strconv.Itoa(row) + "," + strconv.Itoa(column))
	})

	table.Select(0, 0).SetFixed(fixPos[0], fixPos[1]).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}

	}).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
	})
	return nil
}
