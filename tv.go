package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	initView()
	args.setDefault()
	RootCmd := &cobra.Command{
		Use:     "tv {File_Name}",
		Version: "0.5.3",
		Short:   "tv(Table Viewer) for delimited file in terminal",
		Run: func(cmd *cobra.Command, cmdargs []string) {
			if args.Sep == "\\t" {
				args.Sep = "	"
			}
			if len([]rune(args.Sep)) > 0 {
				b.sep = []rune(args.Sep)[0]
			}
			info, err := os.Stdin.Stat()
			fatalError(err)
			//check whether from a console pipe
			if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
				if len(cmdargs) < 1 {
					stopView()
					_ = cmd.Help()
					return
				}
				//get file name form console
				args.FileName = cmdargs[0]
				usefulInfo("Data loading...")
				err = loadFileToBuffer(args.FileName, b)
				fatalError(err)
				usefulInfo("Data loaded")
			} else {
				args.FileName = "From Shell Pipe"
				usefulInfo("Data loading...")
				err = loadPipeToBuffer(os.Stdin, b)
				usefulInfo("Data loaded")
				fatalError(err)
			}

			//process freeze mode
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
			err = drawUI(b, args.Transpose)
			fatalError(err)
			if !debug {
				if err = app.SetRoot(UI, true).SetFocus(UI).Run(); err != nil {
					fatalError(err)
				}
			}
		},
	}

	RootCmd.Flags().StringVar(&args.Sep, "s", "", "(optional) Split symbol")
	RootCmd.Flags().IntVar(&args.NLine, "nl", 0, "(optional) Only display first N lines")
	RootCmd.Flags().StringSliceVar(&args.SkipSymbol, "is", []string{}, "(optional) Ignore lines with specific prefix(multiple arguments support, separated by comma)")
	RootCmd.Flags().IntVar(&args.SkipNum, "in", 0, "(optional) Ignore first N row")
	RootCmd.Flags().IntSliceVar(&args.ShowNum, "dc", []int{}, "(optional) Only display certain columns(multiple arguments support, separated by comma)")
	RootCmd.Flags().IntSliceVar(&args.HideNum, "hc", []int{}, "(optional) Do not display certain columns(multiple arguments support, separated by comma)")
	RootCmd.Flags().IntVar(&args.Header, "fi", 0, "(optional) [default: 0]\n-1, Unfreeze first row and first column\n 0, Freeze first row and first column\n 1, Freeze first row\n 2, Freeze first column")
	RootCmd.Flags().BoolVar(&args.Transpose, "tr", false, "(optional) Transpose data")
	RootCmd.Flags().BoolVar(&args.Strict, "strict", false, "(optional) Check for missing data")
	RootCmd.Flags().SortFlags = false
	err := RootCmd.Execute()
	fatalError(err)
}
