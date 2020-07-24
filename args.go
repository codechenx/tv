package main

//Args struct
type Args struct {
	FileName   string
	Sep        string
	SkipSymbol []string //ignore line with specified prefix
	SkipNum    int      //Number of lines that should be skipped
	ShowNum    []int    //columns that should be displayed
	HideNum    []int    //columns that should be hidden
	Header     int      //header display mode
	Transpose  bool
	NLine      int //Number of lines that should be displayed
}

func (args *Args) setDefault() {
	args.Sep = ""
	args.SkipSymbol = []string{}
	args.SkipNum = 0
	args.ShowNum = []int{}
	args.HideNum = []int{}
	args.Header = 0
	args.Transpose = false
	args.NLine = 0
}
