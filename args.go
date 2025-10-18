package main

// Args struct
type Args struct {
	FileName   string
	Sep        string
	SkipSymbol []string //ignore line with specified prefix
	SkipNum    int      //Number of lines that should be skipped
	ShowNum    []int    //columns that should be displayed
	HideNum    []int    //columns that should be hidden
	Header     int      //header display mode
	NLine      int      //number of lines that should be displayed
	Strict     bool     // check for missing data
	AsyncLoad  bool     // enable async loading for progressive rendering
}

func (args *Args) setDefault() {
	args.Sep = ""
	args.SkipSymbol = []string{}
	args.SkipNum = 0
	args.ShowNum = []int{}
	args.HideNum = []int{}
	args.Header = 0
	args.NLine = 0
	args.Strict = false
	args.AsyncLoad = true // default to async loading
}
