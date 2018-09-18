package main

import (
	"bufio"
	"compress/gzip"
	"errors"
	"os"
	"strings"
)

func loadFile(fn string, b *Buffer) error {
	if !exists(fn) {
		return errors.New("the file does not exist")
	}
	comp := compressed(fn)
	scanner, err := fScanner(fn, comp)
	if err != nil {
		return err
	}
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
			b.sep, err = deterSep(line)
		}
		if err != nil {
			return err
		}

		if len(args.ShowNum) != 0 || len(args.HideNum) != 0 {
			var lineSli []string
			tempLineSli := strings.Split(line, b.sep)
			visCol, err := getVisCol(args.ShowNum, args.HideNum, len(tempLineSli))
			if err != nil {
				return err
			}
			for _, i := range visCol {
				lineSli = append(lineSli, tempLineSli[i])
			}
			err = b.contAppendSli(lineSli, true)
			if err != nil {
				return err
			}

		} else {
			err = b.contAppendStr(line, b.sep, true)
			if err != nil {
				return err
			}
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
	for _, sy := range sy {
		if strings.HasPrefix(line, sy) {
			return true
		}

	}
	return false
}

func compressed(fn string) bool {
	if strings.HasSuffix(fn, ".gz") {
		return true
	}
	return false
}

func deterSep(line string) (string, error) {
	if strings.Contains(line, "\t") {
		return "\t", nil
	} else if strings.Contains(line, ",") {
		return ",", nil
	}
	return "", errors.New("tv can't identify separator, you need to set it manual")
}

func fScanner(fn string, comp bool) (*bufio.Scanner, error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	if comp {
		gzCont, err := gzip.NewReader(file)
		if err != nil {
			return nil, err
		}
		return bufio.NewScanner(gzCont), nil
	}
	return bufio.NewScanner(file), nil
}

func getVisCol(showNumL, hideNumL []int, colLen int) ([]int, error) {
	for _, i := range showNumL {
		if i > colLen {
			return nil, errors.New("column " + string(i) + "does not exist")
		}
	}

	for _, i := range hideNumL {
		if i > colLen {
			return nil, errors.New("column " + string(i) + "does not exist")
		}
	}

	var visCol []int
	for i := 0; i < colLen; i++ {
		flag, err := checkVisible(showNumL, hideNumL, i)
		if err != nil {
			return nil, err
		}
		if flag {
			visCol = append(visCol, i)
		}
	}
	return visCol, nil

}

func checkVisible(showNumL, hideNumL []int, col int) (bool, error) {
	if len(showNumL) != 0 && len(hideNumL) != 0 {
		return false, errors.New("you can only set visible column or hidden column")
	}

	if len(showNumL) != 0 {
		for _, colTestS := range showNumL {
			if col+1 == colTestS {
				return true, nil
			}
		}
		return false, nil
	}
	if len(hideNumL) != 0 {
		for _, colTestH := range hideNumL {
			if col+1 == colTestH {
				return false, nil
			}
		}
	}
	return true, nil
}
