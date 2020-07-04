package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"strconv"
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

//print useful info and force quite app
func usefulInfo(s string) {
	color.Set(color.FgHiYellow)
	fmt.Println(s)
	color.Unset()
}

//I2B  covert int to bool, if i >0:true, else false
func I2B(i int) bool {
	if i > 0 {
		return true
	}
	return false
}

//F2S covert float64 to bool
func F2S(i float64) string {
	return strconv.FormatFloat(i, 'f', 4, 64)
}

//S2F covert string to float64
func S2F(i string) float64 {
	s, err := strconv.ParseFloat(i, 64)
	if err != nil {
		fatalError(err)
	}
	return s
}

//I2S covert int to string
func I2S(i int) string {
	return strconv.Itoa(i)
}

func getHelpContent() string {
	helpContent := `
C == Ctrl

##Quit##
ESC                 Quit

##Movement##
Left-arrow          Move left
Right-arrow         Move right
Down-arrow          Move down
UP-arrow            Move up

h                   Move left
l                   Move right
j                   Move down
k                   Move up

C-F                 Move down by one page 
C-B                 Move up by one page  

C-e                 Move to end of current column
C-h                 Move to head of current column

G                   Move to last cell of table
g                   Move to first cell of table

##Data Type##
C-m 				Change column data type to string or number

##Sort##
C-k                 Sort data by column(ascend)
C-l                 Sort data by column(descend)

##Stats##
C-y                 Show basic stats of current column, back to data table
`
	return helpContent
}
