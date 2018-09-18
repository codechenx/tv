package main

import (
	"fmt"
	"github.com/fatih/color"
)

func fatalError(err error) {
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println(err)
		color.Unset()
		if app != nil {
			app.Stop()
		}
	}
}

func warningError(err error) {
	if err != nil {
		color.Set(color.FgYellow)
		fmt.Println(err)
		color.Unset()
		if app != nil {
			app.Stop()
		}
	}
}
