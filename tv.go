package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	initView()
	args.setDefault()
	RootCmd := &cobra.Command{
		Use:   "tv {File_Name}",
		Short: "tv(Table Viewer) for delimited file in terminal",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, cmdargs []string) {

			var err error
			args.FileName = cmdargs[0]
			if args.Sep == "\\t" {
				b.sep = "\t"
			} else {
				b.sep = args.Sep
			}

			switch args.Header {
			case -1:
				b.rowFreeze, b.colFreeze = 0, 0
			case 0:
				b.rowFreeze, b.colFreeze = 1, 1
			case 1:
				b.rowFreeze, b.colFreeze = 1, 0
			case 2:
				b.rowFreeze, b.colFreeze = 0, 1

			}

			err = loadFile(args.FileName, b)
			fatalError(err)
			err = drawUI(b, args.Transpose)
			fatalError(err)
			if !debug {
				if err = app.SetRoot(UI, true).SetFocus(UI).Run(); err != nil {
					panic(err)
				}
			}
		},
	}

	RootCmd.Flags().StringVar(&args.Sep, "s", "", "split symbol [default: \"\"]")
	RootCmd.Flags().StringSliceVar(&args.SkipSymbol, "is", []string{}, "ignore lines with specific prefix(support for multiple arguments, separated by comma")
	RootCmd.Flags().IntVar(&args.SkipNum, "ir", 0, "ignore first N row [default: 0]")
	RootCmd.Flags().IntSliceVar(&args.ShowNum, "dc", []int{}, "only display certain columns(support for multiple arguments, separated by comma)")
	RootCmd.Flags().IntSliceVar(&args.HideNum, "hc", []int{}, "do not display certain columns(support for multiple arguments, separated by comma)")
	RootCmd.Flags().IntVar(&args.Header, "fi", 0, "-1, Unfreeze first row and first column; 0, Freeze first row and first column; 1, Freeze first row; 2, Freeze first column [default: 0]")
	RootCmd.Flags().BoolVar(&args.Transpose, "tr", false, "transpose and view data [default: false]")
	RootCmd.Flags().SortFlags = false
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
