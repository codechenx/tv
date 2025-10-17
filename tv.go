package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func main() {
	initView()
	args.setDefault()
	RootCmd := &cobra.Command{
		Use:     "tv {File_Name}",
		Version: "0.6.1",
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

			// Determine if we should use async loading
			useAsync := args.AsyncLoad

			//check whether from a console pipe
			if info.Mode()&os.ModeCharDevice != 0 {
				if len(cmdargs) < 1 {
					stopView()
					_ = cmd.Help()
					return
				}
				//get file name form console
				args.FileName = cmdargs[0]

				// Check if file exists before attempting to load
				if _, err := os.Stat(args.FileName); os.IsNotExist(err) {
					stopView()
					fmt.Printf("⚠️  File not found: %s\n", args.FileName)
					os.Exit(1)
				} else if err != nil {
					stopView()
					fmt.Printf("⚠️  Cannot access file: %s\n", err)
					os.Exit(1)
				}

				if useAsync {
					// Start async loading
					userMovedCursor = false // Reset cursor tracking
					updateChan := make(chan bool, 10)
					doneChan := make(chan error, 1)
					go loadFileToBufferAsync(args.FileName, b, updateChan, doneChan)

					// Wait for initial data or error
					select {
					case <-updateChan:
						// Initial data ready
					case err := <-doneChan:
						// Error during initial loading
						fatalError(err)
						return
					}

					// Process freeze mode
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

					// Check if file is empty (no data rows)
					dataRows := b.rowLen - b.rowFreeze
					if b.rowLen == 0 || dataRows <= 0 {
						stopView()
						if b.rowLen == 0 {
							fmt.Println("⚠️  File is empty (no rows)")
						} else {
							fmt.Println("⚠️  File is empty (only header, no data rows)")
						}
						os.Exit(0)
					}

					// Draw initial UI
					err = drawUI(b, args.Transpose)
					fatalError(err)

					// Start update handler in background
					go func() {
						ticker := time.NewTicker(20 * time.Millisecond)
						defer ticker.Stop()

						loadComplete := false
						for !loadComplete {
							select {
							case <-updateChan:
								// Update available - will be handled by ticker
							case err := <-doneChan:
								loadComplete = true
								if err != nil {
									fatalError(err)
								}
								// Final update
								app.QueueUpdateDraw(func() {
									drawBuffer(b, bufferTable, args.Transpose)
									updateFooterWithStatus("Loaded " + strconv.Itoa(b.rowLen) + " rows")
								})
							case <-ticker.C:
								// Periodic UI update
								app.QueueUpdateDraw(func() {
									drawBuffer(b, bufferTable, args.Transpose)

									// Keep cursor on first row if user hasn't moved it
									if !userMovedCursor {
										row, col := bufferTable.GetSelection()
										if row != 0 {
											bufferTable.Select(0, col)
										}
									}

									if loadProgress.TotalBytes > 0 {
										// Show percentage for files
										percent := loadProgress.GetPercentage()
										updateFooterWithStatus(fmt.Sprintf("Loading... %.1f%%", percent))
									} else {
										// Show row count for pipes (no file size)
										updateFooterWithStatus("Loading... " + strconv.Itoa(b.rowLen) + " rows")
									}
								})
							}
						}
					}()

					if !debug {
						if err = app.SetRoot(UI, true).SetFocus(UI).Run(); err != nil {
							fatalError(err)
						}
					}
				} else {
					// Synchronous loading (original behavior)
					err = loadFileToBuffer(args.FileName, b)
					fatalError(err)

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

					// Check if file is empty (no data rows)
					dataRows := b.rowLen - b.rowFreeze
					if b.rowLen == 0 || dataRows <= 0 {
						stopView()
						if b.rowLen == 0 {
							fmt.Println("⚠️  File is empty (no rows)")
						} else {
							fmt.Println("⚠️  File is empty (only header, no data rows)")
						}
						os.Exit(0)
					}
					err = drawUI(b, args.Transpose)
					fatalError(err)
					if !debug {
						if err = app.SetRoot(UI, true).SetFocus(UI).Run(); err != nil {
							fatalError(err)
						}
					}
				}
			} else {
				args.FileName = "From Shell Pipe"

				if useAsync {
					// Start async loading
					userMovedCursor = false // Reset cursor tracking
					updateChan := make(chan bool, 10)
					doneChan := make(chan error, 1)
					go loadPipeToBufferAsync(os.Stdin, b, updateChan, doneChan)

					// Wait for initial data or error
					select {
					case <-updateChan:
						// Initial data ready
					case err := <-doneChan:
						// Error during initial loading
						fatalError(err)
						return
					}

					// Process freeze mode
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

					// Check if pipe is empty (no data rows)
					dataRows := b.rowLen - b.rowFreeze
					if b.rowLen == 0 || dataRows <= 0 {
						stopView()
						if b.rowLen == 0 {
							fmt.Println("⚠️  No data received from pipe (empty input)")
						} else {
							fmt.Println("⚠️  No data received from pipe (only header, no data rows)")
						}
						os.Exit(0)
					}

					// Draw initial UI
					err = drawUI(b, args.Transpose)
					fatalError(err)

					// Start update handler in background
					go func() {
						ticker := time.NewTicker(20 * time.Millisecond)
						defer ticker.Stop()

						loadComplete := false
						for !loadComplete {
							select {
							case <-updateChan:
								// Update available - will be handled by ticker
							case err := <-doneChan:
								loadComplete = true
								if err != nil {
									fatalError(err)
								}
								// Final update
								app.QueueUpdateDraw(func() {
									drawBuffer(b, bufferTable, args.Transpose)
									updateFooterWithStatus("Loaded " + strconv.Itoa(b.rowLen) + " rows")
								})
							case <-ticker.C:
								// Periodic UI update
								app.QueueUpdateDraw(func() {
									drawBuffer(b, bufferTable, args.Transpose)

									// Keep cursor on first row if user hasn't moved it
									if !userMovedCursor {
										row, col := bufferTable.GetSelection()
										if row != 0 {
											bufferTable.Select(0, col)
										}
									}

									if loadProgress.TotalBytes > 0 {
										// Show percentage for files
										percent := loadProgress.GetPercentage()
										updateFooterWithStatus(fmt.Sprintf("Loading... %.1f%%", percent))
									} else {
										// Show row count for pipes (no file size)
										updateFooterWithStatus("Loading... " + strconv.Itoa(b.rowLen) + " rows")
									}
								})
							}
						}
					}()

					if !debug {
						if err = app.SetRoot(UI, true).SetFocus(UI).Run(); err != nil {
							fatalError(err)
						}
					}
				} else {
					// Synchronous loading (original behavior)
					err = loadPipeToBuffer(os.Stdin, b)
					fatalError(err)

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

					// Check if pipe is empty (no data rows)
					dataRows := b.rowLen - b.rowFreeze
					if b.rowLen == 0 || dataRows <= 0 {
						stopView()
						if b.rowLen == 0 {
							fmt.Println("⚠️  No data received from pipe (empty input)")
						} else {
							fmt.Println("⚠️  No data received from pipe (only header, no data rows)")
						}
						os.Exit(0)
					}
					err = drawUI(b, args.Transpose)
					fatalError(err)
					if !debug {
						if err = app.SetRoot(UI, true).SetFocus(UI).Run(); err != nil {
							fatalError(err)
						}
					}
				}
			}
		},
	}

	RootCmd.Flags().StringVarP(&args.Sep, "separator", "s", "", "Delimiter/separator character (use \\t for tab)")
	RootCmd.Flags().IntVarP(&args.NLine, "lines", "n", 0, "Display only first N lines")
	RootCmd.Flags().StringSliceVar(&args.SkipSymbol, "skip-prefix", []string{}, "Skip lines starting with prefix (comma-separated)")
	RootCmd.Flags().IntVar(&args.SkipNum, "skip-lines", 0, "Skip first N lines")
	RootCmd.Flags().IntSliceVar(&args.ShowNum, "columns", []int{}, "Show only specified columns (comma-separated)")
	RootCmd.Flags().IntSliceVar(&args.HideNum, "hide-columns", []int{}, "Hide specified columns (comma-separated)")
	RootCmd.Flags().IntVarP(&args.Header, "freeze", "f", 0, "Freeze mode: -1=none, 0=row+col, 1=row only, 2=col only")
	RootCmd.Flags().BoolVarP(&args.Transpose, "transpose", "t", false, "Transpose rows and columns")
	RootCmd.Flags().BoolVar(&args.Strict, "strict", false, "Strict mode: fail on missing/inconsistent data")
	RootCmd.Flags().BoolVar(&args.AsyncLoad, "async", true, "Progressive rendering while loading")
	RootCmd.Flags().SortFlags = false
	err := RootCmd.Execute()
	fatalError(err)
}
