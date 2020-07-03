package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

//print fatal error and force quite app
func fatalError(err error) {
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println(err)
		color.Unset()
		if app != nil {
			app.Stop()
		}
		if !debug {
			os.Exit(1)
		}
	}
}
