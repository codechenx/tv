package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func render(b *Buffer) {
	cols, rows := b.colNum+b.vHCN, b.rowNum+b.vHRN

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			if c < 1 || r < 1 {
				color = tcell.ColorYellow
			}
			if !args.Transpose {
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

	table.SetFixed(1, 1)
}
