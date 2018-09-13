package main

import (
	"github.com/alexflint/go-arg"
	"strings"
)

func main() {
	args.setDefault()
	arg.MustParse(&args) // temp
	initView()
	b := createNewBuffer()
	b.header = args.Header
	b.sep = args.Sep
	comp := false
	if strings.HasSuffix(args.FileName, ".gz") {
		comp = true
	}
	err := loadFile(args.FileName, b, comp)
	fatalError(err)
	b.addVirHeader()
	render(b)
	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}
