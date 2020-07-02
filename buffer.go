package main

import (
	"errors"
	"sort"
	"strings"
)

//buffer to store data
type Buffer struct {
	sep       string
	cont      [][]string
	rowLen    int
	colLen    int
	rowFreeze int // true :1, false:0
	colFreeze int // true :1, false:0
}

func createNewBuffer() *Buffer {
	return &Buffer{sep: "", cont: [][]string{}, rowLen: 0, colLen: 0, rowFreeze: 1, colFreeze: 1}
}

func (b *Buffer) contAppendSli(s []string, strict bool) error {
	if b.rowLen == 0 {
		b.colLen = len(s)
	}
	if strict && len(s) != b.colLen {
		return errors.New("lack some column")
	}

	b.cont = append(b.cont, s)
	b.rowLen++
	return nil
}

func (b *Buffer) contAppendStr(s, sep string, strict bool) error {
	sSli := strings.Split(s, sep)
	if b.rowLen == 0 {
		b.colLen = len(sSli)
	}

	if strict && len(sSli) != b.colLen {
		return errors.New("lack some column")
	}

	b.cont = append(b.cont, sSli)
	b.rowLen++
	return nil
}

func (b *Buffer) sort(colIndex int, rev bool) {
	if rev {
		if I2B(b.rowFreeze) {
			sort.SliceStable(b.cont[1:], func(i, j int) bool { return b.cont[1:][i][colIndex] > b.cont[1:][j][colIndex] })
		} else {
			sort.SliceStable(b.cont, func(i, j int) bool { return b.cont[i][colIndex] > b.cont[j][colIndex] })
		}
	} else {

		if I2B(b.rowFreeze) {
			sort.SliceStable(b.cont[1:], func(i, j int) bool { return b.cont[1:][i][colIndex] < b.cont[1:][j][colIndex] })
		} else {
			sort.SliceStable(b.cont, func(i, j int) bool { return b.cont[i][colIndex] < b.cont[j][colIndex] })
		}
	}
}
