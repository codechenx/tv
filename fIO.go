package main

import (
	"bufio"
	"compress/gzip"
	"errors"
	"os"
	"strings"
)

func loadFile(fn string, b *Buffer, comp bool) error {
	if !exists(fn) {
		return errors.New("the file does not exist")
	}
	file, err := os.Open(fn)
	if err != nil {
		return err
	}
	var scanner *bufio.Scanner
	if comp {
		gzCont, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		scanner = bufio.NewScanner(gzCont)
	} else {
		scanner = bufio.NewScanner(file)
	}

	//scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		//ignore first n lines
		if args.SkipNum > 0 {
			args.SkipNum--
			continue
		}
		//ignore line with specified prefix
		if skipLine(line, args.SkipSymbol) {
			continue
		}
		//Auto set split symbols
		if b.sep == "" {
			if strings.Contains(line, "\t") {
				b.sep = "\t"
			} else if strings.Contains(line, ",") {
				b.sep = ","
			} else {
				fatalError(errors.New("you need to set a separator"))
			}
		} else if b.sep == "\\t" {
			b.sep = "\t"
		}
		err = b.contAppend(line, b.sep, true)
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

func skipLine(line string, sy []string) bool {
	if len(sy) != 0 {
		for _, sy := range sy {
			if strings.HasPrefix(line, sy) {
				return true
			}
		}
	}
	return false
}
