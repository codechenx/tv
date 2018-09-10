package main

import "github.com/rivo/tview"

var app *tview.Application
var table *tview.Table
var args Args

func initView() {
	app = tview.NewApplication()
	table = tview.NewTable()
	table.SetBorders(true)
}
