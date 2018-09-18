package main

import (
	"errors"
	"strconv"
	"strings"
)

//buffer to store data
type Buffer struct {
	name   string
	sep    string
	header int
	vHRN   int // virHeaderRowNumber
	vHCN   int //virHeaderColNumber
	cont   [][]string
	rowNum int
	colNum int
}

func createNewBuffer() *Buffer {
	return &Buffer{name: "", sep: "", header: 0, vHRN: 0, vHCN: 0, rowNum: 0, colNum: 0, cont: [][]string{}}
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

func (b *Buffer) addVirHeader() {
	var rowVirHeader []string
	var colVirHeader []string
	for i := 0; i < b.colNum; i++ {
		colVirHeader = append(colVirHeader, strconv.Itoa(i+1))
	}
	rowVirHeader = append(rowVirHeader, string("#"))
	for i := 0; i < b.rowNum; i++ {
		rowVirHeader = append(rowVirHeader, strconv.Itoa(i+1))
	}
	if b.header == -1 {
		b.cont = append(b.cont, []string{})
		copy(b.cont[1:], b.cont)
		b.cont[0] = colVirHeader
		for i := 0; i < b.rowNum+1; i++ {
			b.cont[i] = append(b.cont[i], "")
			copy(b.cont[i][1:], b.cont[i])
			b.cont[i][0] = rowVirHeader[i]
		}
		b.vHRN++
		b.vHCN++
	}

	if b.header == 0 {
		for i := 0; i < b.rowNum; i++ {
			b.cont[i] = append(b.cont[i], "")
			copy(b.cont[i][1:], b.cont[i])
			b.cont[i][0] = rowVirHeader[i]
		}
		b.vHCN++
	}

	if b.header == 1 {
		b.cont = append(b.cont, []string{})
		copy(b.cont[1:], b.cont)
		b.cont[0] = colVirHeader
		b.vHRN++
	}

	if b.header == 2 {
		return
	}
}
