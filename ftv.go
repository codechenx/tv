package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// setupFreezeMode configures row and column freeze settings based on header mode
func setupFreezeMode(b *Buffer) {
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
}

// validateDataNotEmpty checks if buffer has data rows and exits if empty
func validateDataNotEmpty(b *Buffer, source string) error {
	dataRows := b.rowLen - b.rowFreeze
	if b.rowLen == 0 || dataRows <= 0 {
		stopView()
		if b.rowLen == 0 {
			fmt.Printf("⚠️  %s is empty (no rows)\n", source)
		} else {
			fmt.Printf("⚠️  %s is empty (only header, no data rows)\n", source)
		}
		os.Exit(0)
	}
	return nil
}

// startAsyncUpdateHandler manages UI updates during async loading
func startAsyncUpdateHandler(updateChan <-chan bool, doneChan <-chan error) {
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
					drawBuffer(b, bufferTable)
					updateFooterWithStatus("Loaded " + strconv.Itoa(b.rowLen) + " rows")
				})
			case <-ticker.C:
				// Periodic UI update
				app.QueueUpdateDraw(func() {
					drawBuffer(b, bufferTable)

					// Keep cursor on first row if user hasn't moved it
					if !userMovedCursor {
						row, col := bufferTable.GetSelection()
						if row != 0 {
							bufferTable.Select(0, col)
						}
					}

					if loadProgress.TotalBytes > 0 {
						// Show progress bar for files
						percent := loadProgress.GetPercentage()
						progressBar := makeProgressBar(percent, 15)
						updateFooterWithStatus(fmt.Sprintf("Loading... %s", progressBar))
					} else {
						// Show row count for pipes (no file size)
						updateFooterWithStatus("Loading... " + strconv.Itoa(b.rowLen) + " rows")
					}
				})
			}
		}
	}()
}

// loadDataAsync starts async loading and waits for initial data
func loadDataAsync(loader func(*Buffer, chan<- bool, chan<- error), b *Buffer) (chan bool, chan error, error) {
	userMovedCursor = false // Reset cursor tracking
	updateChan := make(chan bool, 10)
	doneChan := make(chan error, 1)

	loader(b, updateChan, doneChan)

	// Wait for initial data or error
	select {
	case <-updateChan:
		// Initial data ready
		return updateChan, doneChan, nil
	case err := <-doneChan:
		// Error during initial loading
		return nil, nil, err
	}
}

// runApp starts the UI application if not in debug mode
func runApp() error {
	if !debug {
		if err := app.SetRoot(UI, true).SetFocus(UI).Run(); err != nil {
			return err
		}
	}
	return nil
}

// loadAndDisplayAsync handles the complete async loading workflow
func loadAndDisplayAsync(loader func(*Buffer, chan<- bool, chan<- error), source string) error {
	updateChan, doneChan, err := loadDataAsync(loader, b)
	if err != nil {
		return err
	}

	setupFreezeMode(b)
	if err := validateDataNotEmpty(b, source); err != nil {
		return err
	}

	if err := drawUI(b); err != nil {
		return err
	}

	startAsyncUpdateHandler(updateChan, doneChan)
	return runApp()
}

// loadAndDisplaySync handles the complete sync loading workflow
func loadAndDisplaySync(loader func(*Buffer) error, source string) error {
	if err := loader(b); err != nil {
		return err
	}

	setupFreezeMode(b)
	if err := validateDataNotEmpty(b, source); err != nil {
		return err
	}

	if err := drawUI(b); err != nil {
		return err
	}

	return runApp()
}

func main() {
	initView()
	args.setDefault()
	RootCmd := &cobra.Command{
		Use:     "ftv {File_Name}",
		Version: "0.8",
		Short:   "Fast table viewer for delimited file in terminal",
		Run: func(cmd *cobra.Command, cmdargs []string) {
			if args.Sep == "\\t" {
				args.Sep = "	"
			}
			if len([]rune(args.Sep)) > 0 {
				b.sep = []rune(args.Sep)[0]
			}

			// Configure memory limit
			if args.MemoryMB > 0 {
				b.setMemoryLimit(int64(args.MemoryMB) * 1024 * 1024) // Convert MB to bytes
			}
			// else use default (unlimited - 0)

			info, err := os.Stdin.Stat()
			fatalError(err)

			// Determine if we should use async loading
			useAsync := args.AsyncLoad

			//check whether from a console pipe
			if info.Mode()&os.ModeCharDevice != 0 {
				// FILE MODE
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
					err = loadAndDisplayAsync(func(b *Buffer, updateChan chan<- bool, doneChan chan<- error) {
						go loadFileToBufferAsync(args.FileName, b, updateChan, doneChan)
					}, "File")
					fatalError(err)
				} else {
					err = loadAndDisplaySync(func(b *Buffer) error {
						return loadFileToBuffer(args.FileName, b)
					}, "File")
					fatalError(err)
				}
			} else {
				// PIPE MODE
				args.FileName = "From Shell Pipe"

				if useAsync {
					err = loadAndDisplayAsync(func(b *Buffer, updateChan chan<- bool, doneChan chan<- error) {
						go loadPipeToBufferAsync(os.Stdin, b, updateChan, doneChan)
					}, "Pipe")
					fatalError(err)
				} else {
					err = loadAndDisplaySync(func(b *Buffer) error {
						return loadPipeToBuffer(os.Stdin, b)
					}, "Pipe")
					fatalError(err)
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
	RootCmd.Flags().BoolVar(&args.Strict, "strict", false, "Strict mode: fail on missing/inconsistent data")
	RootCmd.Flags().BoolVar(&args.AsyncLoad, "async", true, "Progressive rendering while loading")
	RootCmd.Flags().IntVarP(&args.MemoryMB, "memory", "m", 0, "Memory limit in MB (0=unlimited/default, >0=set limit)")
	RootCmd.Flags().SortFlags = false
	err := RootCmd.Execute()
	fatalError(err)
}
