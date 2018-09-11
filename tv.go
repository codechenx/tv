package main

import (
	"github.com/alexflint/go-arg"
)

func main() {
	args.setDefault()
	arg.MustParse(&args) // temp
	initView()
	b := createNewBuffer()
	b.header = args.Header
	b.sep = args.Sep
	err := loadFile(args.FileName, b)
	fatalError(err)
	b.addVirHeader()
	render(b)
	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}
