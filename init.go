package main

import (
	"github.com/rivo/tview"
)

var app *tview.Application
var table *tview.Table
var grid *tview.Grid
var args Args
var debug bool

func initView() {
	app = tview.NewApplication()
	table = tview.NewTable()
	table.SetBorders(true)
}

func stopView() {
	app.Stop()
}
