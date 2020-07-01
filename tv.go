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
		Short: "An example cobra command",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, cmdargs []string) {

			var err error
			args.FileName = cmdargs[0]
			b := createNewBuffer()
			b.header = args.Header
			if args.Sep == "\\t" {
				b.sep = "\t"
			} else {
				b.sep = args.Sep
			}

			err = loadFile(args.FileName, b)
			fatalError(err)
			b.addVirHeader()
			err = render(b, args.Transpose)
			fatalError(err)
			if !debug {
				if err = app.SetRoot(table, true).Run(); err != nil {
					panic(err)
				}
			}
		},
	}

	RootCmd.Flags().StringVar(&args.Sep, "s", "", "split symbol [default: \"\"]")
	RootCmd.Flags().StringSliceVar(&args.SkipSymbol, "is", []string{}, "ignore lines with specific prefix(support for multiple arguments, separated by space")
	RootCmd.Flags().IntVar(&args.SkipNum, "ir", 0, "ignore first N row [default: 0]")
	RootCmd.Flags().IntSliceVar(&args.ShowNum, "dc", []int{}, "only display certain columns(support for multiple arguments, separated by space)")
	RootCmd.Flags().IntSliceVar(&args.HideNum, "hc", []int{}, "do not display certain columns(support for multiple arguments, separated by space)")
	RootCmd.Flags().IntVar(&args.Header, "he", 0, "-1, no column name and row name; 0, use first row as row name; 1, use first column as column name; 2, use first column as column name and first row as row name [default: 0]")
	RootCmd.Flags().BoolVar(&args.Transpose, "tr", false, "transpose and view data [default: false]")
	RootCmd.Flags().SortFlags = false
	err := RootCmd.Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
