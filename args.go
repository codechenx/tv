package main

//Args struct
type Args struct {
	FileName   string
	Sep        string
	SkipSymbol []string
	SkipNum    int
	ShowNum    []int
	HideNum    []int
	Header     int
	Transpose  bool
}

func (args *Args) setDefault() {
	args.Sep = ""
	args.SkipSymbol = []string{}
	args.SkipNum = 0
	args.ShowNum = []int{}
	args.HideNum = []int{}
	args.Header = 0
	args.Transpose = false
}

func (Args) Version() string {
	return "tv 0.4"
}
func (Args) Description() string {
	return "tv(Table Viewer) for delimited file in terminal "
}
