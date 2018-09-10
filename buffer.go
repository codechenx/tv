package main

import (
	"errors"
	"strconv"
	"strings"
)

type Buffer struct {
	name   string
	sep    string
	header int
	cont   [][]string
	rowNum int
	colNum int
}

func createNewBuffer() *Buffer {
	return &Buffer{name: "", sep: "", header: 0, rowNum: 0, colNum: 0, cont: [][]string{}}
}

func (b *Buffer) contAppend(s, sep string) error {
	sSlice := strings.Split(s, sep)
	if b.rowNum == 0 {
		b.colNum = len(sSlice)
	}
	if len(sSlice) != b.colNum {
		return errors.New("lack some column")
	}
	b.cont = append(b.cont, sSlice)
	b.rowNum++
	return nil
}

func (b Buffer) addVirHeader() {
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
	}

	if b.header == 0 {
		for i := 0; i < b.rowNum; i++ {
			b.cont[i] = append(b.cont[i], "")
			copy(b.cont[i][1:], b.cont[i])
			b.cont[i][0] = rowVirHeader[i]
		}
	}

	if b.header == 1 {
		b.cont = append(b.cont, []string{})
		copy(b.cont[1:], b.cont)
		b.cont[0] = colVirHeader
	}

	if b.header == 2 {
		return
	}
}
