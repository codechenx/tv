package main

import (
	"errors"
	"strings"
)

//buffer to store data
type Buffer struct {
	sep    string
	header int
	cont   [][]string
	rowNum int
	colNum int
}

func createNewBuffer() *Buffer {
	return &Buffer{sep: "", header: 0, cont: [][]string{}, rowNum: 0, colNum: 0}
}

func (b *Buffer) contAppendSli(s []string, strict bool) error {
	if b.rowNum == 0 {
		b.colNum = len(s)
	}
	if strict && len(s) != b.colNum {
		return errors.New("lack some column")
	}

	b.cont = append(b.cont, s)
	b.rowNum++
	return nil
}

func (b *Buffer) contAppendStr(s, sep string, strict bool) error {
	sSli := strings.Split(s, sep)
	if b.rowNum == 0 {
		b.colNum = len(sSli)
	}

	if strict && len(sSli) != b.colNum {
		return errors.New("lack some column")
	}

	b.cont = append(b.cont, sSli)
	b.rowNum++
	return nil
}
