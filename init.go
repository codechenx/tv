package main

import (
	"github.com/rivo/tview"
)

var app *tview.Application
var UI *tview.Pages
var b *Buffer
var args Args
var debug bool

// initialize tview, buffer, table
func initView() {
	app = tview.NewApplication()
	b = createNewBuffer()
}

//stop UI
func stopView() {
	app.Stop()
}

//covert int to bool, if i >0:true, else false
func I2B(i int) bool {
	if i > 0 {
		return true
	}
	return false
}

// check if input([]string) is digitized
func checkAllNum(a []string) {

}
