package main

import (
	"errors"
	"sort"
)

//buffer struct
type Buffer struct {
	sep       rune
	cont      [][]string
	colType   []int //const colType_str = 0, const colType_int = 0
	rowLen    int
	colLen    int
	rowFreeze int // true:1, false:0
	colFreeze int // true:1, false:0
}

func createNewBuffer() *Buffer {
	return &Buffer{sep: 0, cont: [][]string{}, colType: []int{}, rowLen: 0, colLen: 0, rowFreeze: 1, colFreeze: 1}
}

func createNewBufferWithData(ss [][]string, strict bool) (*Buffer, error) {
	b = createNewBuffer()
	for _, s := range ss {
		err := b.contAppendSli(s, strict)
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

//add []string to buffer
func (b *Buffer) contAppendSli(s []string, strict bool) error {
	if b.rowLen == 0 {
		b.colLen = len(s)
		b.colType = make([]int, b.colLen+1)
	}
	if strict && len(s) != b.colLen {
		return errors.New("Row " + I2S(b.rowLen) + " lack some column")
	}

	b.cont = append(b.cont, s)
	b.rowLen++

	return nil
}

// sort column by string format
func (b *Buffer) sortByStr(colIndex int, rev bool) {
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

// sort column by number format
func (b *Buffer) sortByNum(colIndex int, rev bool) {
	if rev {
		if I2B(b.rowFreeze) {
			sort.SliceStable(b.cont[1:], func(i, j int) bool { return S2F(b.cont[1:][i][colIndex]) > S2F(b.cont[1:][j][colIndex]) })
		} else {
			sort.SliceStable(b.cont, func(i, j int) bool { return S2F(b.cont[i][colIndex]) > S2F(b.cont[j][colIndex]) })
		}
	} else {

		if I2B(b.rowFreeze) {
			sort.SliceStable(b.cont[1:], func(i, j int) bool { return S2F(b.cont[1:][i][colIndex]) < S2F(b.cont[1:][j][colIndex]) })
		} else {
			sort.SliceStable(b.cont, func(i, j int) bool { return S2F(b.cont[i][colIndex]) < S2F(b.cont[j][colIndex]) })
		}
	}
}

// transpose buffer content
func (b *Buffer) transpose() {
	b.rowLen, b.colLen = b.colLen, b.rowLen
	b.colType = make([]int, b.colLen+1)
	xl := len(b.cont[0])
	yl := len(b.cont)
	result := make([][]string, xl)
	for i := range result {
		result[i] = make([]string, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = b.cont[j][i]
		}
	}
	b.cont = result
}

//get ith column []string data
func (b Buffer) getCol(i int) []string {
	result := make([]string, b.rowLen)
	for rowI := 0; rowI < b.rowLen; rowI++ {
		result[rowI] = b.cont[rowI][i]
	}
	return result
}

//set ith column data type
func (b *Buffer) setColType(i int, t int) {
	b.colType[i] = t
}

//get ith column data type
func (b *Buffer) getColType(i int) int {
	return b.colType[i]
}
