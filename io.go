package main

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
)

//load file content to buffer
func loadFileToBuffer(fn string, b *Buffer) error {
	totalAddedLN := 0 //the number of lines has been added into buffer
	scanner, err := getFileScanner(fn)
	if err != nil {
		return err
	}
	scanner.Split(bufio.ScanLines)
	//set separator, if user does not provide it.
	var detectLines []string //lines as detect separator data
	if b.sep == 0 {
		//read 10 lines to detect separator
		lineNumber := 10
		for scanner.Scan() {
			line := scanner.Text()
			//skip empty line
			if line == "\n" {
				continue
			}
			//ignore first n lines
			if args.SkipNum > 0 {
				args.SkipNum--
				continue
			}
			//ignore line with specified prefix
			if skipLine(line, args.SkipSymbol) {
				continue
			}
			detectLines = append(detectLines, line)
			if len(detectLines) >= lineNumber {
				break
			}
		}
		//if the suffix of file name is ".csv", set separator to ",".
		//if the suffix of file name is "tsv", set separator to "\t".
		if strings.HasSuffix(fn, ".csv") {
			b.sep = ','
		} else if strings.HasSuffix(fn, ".tsv") {
			b.sep = '\t'
		} else {
			sd := sepDetecor{}
			b.sep = sd.sepDetect(detectLines)
		}

	}
	//check final separator
	if b.sep == 0 {
		fatalError(errors.New("tv can't identify separator, you need to set it manual"))
	}

	//add detectLines to buffer
	for _, line := range detectLines {
		//parse and add line to buffer
		err = addDRToBuffer(b, line, args.ShowNum, args.HideNum)
		if err != nil {
			return err

		}
		totalAddedLN++
		if totalAddedLN >= args.NLine && args.NLine > 0 {
			break
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		//skip empty line
		if line == "\n" {
			continue
		}
		//ignore first n lines
		if args.SkipNum > 0 && args.NLine > 0 {
			args.SkipNum--
			continue
		}
		//ignore line with specified prefix
		if skipLine(line, args.SkipSymbol) {
			continue
		}

		//parse and add line to buffer
		if totalAddedLN >= args.NLine && args.NLine > 0 {
			break
		}
		err = addDRToBuffer(b, line, args.ShowNum, args.HideNum)
		if err != nil {
			return err
		}
		totalAddedLN++
	}

	return nil
}

//load console pipe content to buffer
func loadPipeToBuffer(stdin io.Reader, b *Buffer) error {
	totalAddedLN := 0 //the number of lines has been added into buffer
	var err error
	scanner := bufio.NewScanner(stdin)
	//read 10 lines to detect separator
	lineNumber := 10
	var detectLines []string //lines as detect separator data
	if b.sep == 0 {
		for scanner.Scan() {
			line := scanner.Text()
			//skip empty line
			if line == "\n" {
				continue
			}
			//ignore first n lines
			if args.SkipNum > 0 {
				args.SkipNum--
				continue
			}
			//ignore line with specified prefix
			if skipLine(line, args.SkipSymbol) {
				continue
			}
			detectLines = append(detectLines, line)
			if len(detectLines) >= lineNumber {
				break
			}
		}
		sd := sepDetecor{}
		b.sep = sd.sepDetect(detectLines)
	}
	//check final separator
	if b.sep == 0 {
		fatalError(errors.New("tv can't identify separator, you need to set it manual"))
	}
	//add detectLines to buffer
	for _, line := range detectLines {
		//parse and add line to buffer
		err = addDRToBuffer(b, line, args.ShowNum, args.HideNum)
		if err != nil {
			return err
		}
		totalAddedLN++
		if totalAddedLN >= args.NLine && args.NLine > 0 {
			break
		}
	}
	for scanner.Scan() {
		line := scanner.Text()
		//skip empty line
		if line == "\n" {
			continue
		}
		//ignore first n lines
		if args.SkipNum > 0 {
			args.SkipNum--
			continue
		}
		//ignore line with specified prefix
		if skipLine(line, args.SkipSymbol) {
			continue
		}

		//parse and add line to buffer
		if totalAddedLN >= args.NLine && args.NLine > 0 {
			break
		}
		err = addDRToBuffer(b, line, args.ShowNum, args.HideNum)
		if err != nil {
			return err
		}
		totalAddedLN++
	}

	return nil
}

//check a line whether should bu skip, according to prefix
func skipLine(line string, sy []string) bool {
	for _, sy := range sy {
		if strings.HasPrefix(line, sy) {
			return true
		}

	}
	return false
}

//get suitable scanner(compressed or not)
func getFileScanner(fn string) (*bufio.Scanner, error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}

	//if input is a gzip file
	if strings.HasSuffix(fn, ".gz") {
		gzCont, err := gzip.NewReader(file)
		if err != nil {
			return nil, err
		}
		return bufio.NewScanner(gzCont), nil
	}

	return bufio.NewScanner(file), nil
}

//check columns that should be displayed
func getVisCol(showNumL, hideNumL []int, colLen int) ([]int, error) {
	for _, i := range showNumL {
		if i > colLen || i <= 0 {
			return nil, errors.New("Column number " + I2S(i) + " does not exist")
		}
	}

	for _, i := range hideNumL {
		if i > colLen || i <= 0 {
			return nil, errors.New("Column number " + I2S(i) + " does not exist")
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

//check ith column should be displayed or not
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

//use go csv library to parse a string line into csv format
func lineCSVParse(s string, sep rune) ([]string, error) {
	r := csv.NewReader(strings.NewReader(s))
	r.Comma = sep
	r.LazyQuotes = true
	//r.TrimLeadingSpace = true //disable, because it will remove NULL item and cause issue.
	record, err := r.Read()
	return record, err
}

//add displayable(according to user's input argument) RowArray(covert line to array) To Buffer
func addDRToBuffer(b *Buffer, line string, showNum, hideNum []int) error {
	var err error
	lineCSVParts, err := lineCSVParse(line, b.sep)
	if err != nil {
		return err
	}
	if len(showNum) != 0 || len(hideNum) != 0 {
		var lineSli []string
		visCol, err := getVisCol(showNum, hideNum, len(lineCSVParts))
		if err != nil {
			return err
		}
		for _, i := range visCol {
			lineSli = append(lineSli, lineCSVParts[i])
		}
		err = b.contAppendSli(lineSli, true)
		if err != nil {
			return err
		}

	} else {
		err := b.contAppendSli(lineCSVParts, true)
		if err != nil {
			return err
		}
	}
	return err
}
