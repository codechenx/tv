package main

import (
	"github.com/rivo/tview"
)

//column data type
const colTypeStr = 0
const colTypeFloat = 1

//get column data type name. s: string, n: number
func type2name(i int) string {
	if i == 0 {
		return "s"
	}
	return "n"
}

var app *tview.Application
var UI *tview.Pages
var b *Buffer
var args Args
var debug bool

// initialize tview, buffer
func initView() {
	app = tview.NewApplication()
	b = createNewBuffer()
}

//stop UI
func stopView() {
	app.Stop()
}
