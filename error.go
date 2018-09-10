package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

func fatalError(err error) {
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println(err)
		color.Unset()
		defer app.Stop()
		os.Exit(1)
	}
}
