package main

import (
	"github.com/alexflint/go-arg"
)

func main() {
	var err error
	args.setDefault()
	arg.MustParse(&args) // temp
	initView()
	b := createNewBuffer()
	b.header = args.Header
	if args.Sep == "\\t" {
		b.sep = "\t"
	} else {
		b.sep = args.Sep
	}
	err = loadFile(args.FileName, b)
	fatalError(err)
	b.addVirHeader()
	err = render(b, args.Transpose)
	fatalError(err)
	if !debug {
		if err = app.SetRoot(table, true).Run(); err != nil {
			panic(err)
		}
	}
}
