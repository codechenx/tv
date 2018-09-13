package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func loadFile(fn string, b *Buffer) error {
	if !exists(fn) {
		return errors.New("file is not exist")
	}
	file, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		//ignore first n lines
		if args.SkipNum > 0 {
			args.SkipNum--
			continue
		}
		//ignore line with specified prefix
		if args.SkipSymbol != "" && strings.HasPrefix(line, args.SkipSymbol) {
			continue
		}
		//Auto set split symbols
		if b.sep == "" {
			if strings.Contains(line, "\t") {
				b.sep = "\t"
			} else if strings.Contains(line, ",") {
				b.sep = ","
			} else {
				fatalError(errors.New("you must set a separator"))
			}
		} else if b.sep == "\\t" {
			b.sep = "\t"
		}
		err = b.contAppend(line, b.sep, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
