package main

type Args struct {
	FileName   string   `arg:"positional"`
	Sep        string   `arg:"-s" help:"split symbol"`
	SkipSymbol []string `arg:"--ss" help:"ignore lines with specific prefix"`
	SkipNum    int      `arg:"--sn" help:"ignore first n lines"`
	Header     int      `arg:"--h" help:" -1, no column name and row name; 0, use first row as row name; 1, use first column as column name; 2, use firt column as column name and first row as row name"`
	Transpose  bool     `arg:"--t" help:"transpose and view data "`
}

func (args *Args) setDefault() {
	args.Sep = ""
	args.SkipSymbol = []string{}
	args.SkipNum = 0
	args.Header = 0
	args.Transpose = false
}

func (Args) Version() string {
	return "tv 0.3"
}
func (Args) Description() string {
	return "tv(Table Viewer) for delimited file in terminal "
}
