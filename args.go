package main

//Args struct
type Args struct {
	FileName   string   `arg:"positional"`
	Sep        string   `arg:"-s" help:"split symbol [default: \"\"]"`
	SkipSymbol []string `arg:"--ss" help:"ignore lines with specific prefix(support for multiple arguments, separated by space)"`
	SkipNum    int      `arg:"--sn" help:"ignore first n lines [default: 0]"`
	ShowNum    []int    `arg:"--rc" help:"show columns(support for multiple arguments, separated by space)"`
	HideNum    []int    `arg:"--hc" help:"hide columns(support for multiple arguments, separated by space)"`
	Header     int      `arg:"--h" help:"-1, no column name and row name; 0, use first row as row name; 1, use first column as column name; 2, use firt column as column name and first row as row name [default: 0]"`
	Transpose  bool     `arg:"--t" help:"transpose and view data [default: false]"`
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
