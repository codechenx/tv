package main

type Args struct {
	FileName   string `arg:"positional"`
	Sep        string `arg:"-s" help:"split symbol"`
	SkipSymbol string `arg:"--ss" help:"skip line with prefix of symbol"`
	Header     int    `arg:"--h" help:" -1, no column name and row name; 0, use first row as row name; 1, use first column as column name; 2, use firt column as column name and first row as row name"`
	Transpose  bool   `arg:"--t" help:"Transpose data and view"`
}

func (args *Args) setDefault() {
	args.Sep = ""
	args.SkipSymbol = ""
	args.Header = 0
	args.Transpose = false
}

func (Args) Version() string {
	return "tv 0.1.1"
}
func (Args) Description() string {
	return "Table Viewer"
}
