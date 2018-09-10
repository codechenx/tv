package main

type Args struct {
	FileName   string `arg:"positional"`
	Sep        string `arg:"-s" help:"split symbol"`
	SkipSymbol string `arg:"--ss" help:"skip line with prefix of symbol"`
	Header     int    `arg:"--h" help:"0 for only ColumnName, 1 for only RowName, 2 for both of ColumnName and RowName"`
}

func (args *Args) setDefault() {
	args.Sep = ""
	args.SkipSymbol = ""
	args.Header = 0
}

func (Args) Version() string {
	return "tv 0.1"
}
func (Args) Description() string {
	return "Table Viewer"
}
