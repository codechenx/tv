package main

import (
	"reflect"
	"testing"
)

func Test_createNewBuffer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createNewBuffer()
		})
	}
}

func TestBuffer_contAppendSli(t *testing.T) {
	type args struct {
		s      []string
		strict bool
	}
	b := createNewBuffer()
	b.colLen = 2
	b.rowLen = 1
	b.cont = [][]string{{"a", "b"}}
	tests := []struct {
		name    string
		b       *Buffer
		args    args
		wantErr bool
	}{
		{"test", b, args{s: []string{"a", "1"}, strict: true}, false},
		{"test", b, args{s: []string{"a", "1", "3"}, strict: true}, true},
		{"test", b, args{s: []string{"a", "1", "2"}, strict: false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.contAppendSli(tt.args.s, tt.args.strict); (err != nil) != tt.wantErr {
				t.Errorf("Buffer.contAppendSli() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createNewBufferWithData(t *testing.T) {
	type args struct {
		ss     [][]string
		strict bool
	}
	wantBuffer := createNewBuffer()
	wantBuffer.colLen = 3
	wantBuffer.rowLen = 4
	wantBuffer.cont = [][]string{{"a", "b", "c"}, {"1", "2", "3"}, {"4", "5", "6"}, {"7", "8", "9"}}
	wantBuffer.colType = []int{0, 0, 0, 0}
	tests := []struct {
		name    string
		args    args
		want    *Buffer
		wantErr bool
	}{
		{"test", args{ss: [][]string{{"a", "b", "c"}, {"1", "2", "3"}, {"4", "5", "6"}, {"7", "8", "9"}}, strict: true}, wantBuffer, false},
		{"test", args{ss: [][]string{{"a", "b", "c"}, {"1", "2", "3"}, {"4", "5"}, {"7", "8", "9"}}, strict: true}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createNewBufferWithData(tt.args.ss, tt.args.strict)
			if (err != nil) != tt.wantErr {
				t.Errorf("createNewBufferWithData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createNewBufferWithData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuffer_sortByStr(t *testing.T) {
	type args struct {
		colIndex int
		rev      bool
	}
	testBuffer, _ := createNewBufferWithData([][]string{{"a", "b", "c"}, {"2", "2", "3"}, {"4", "5", "6"}, {"10", "8", "9"}}, true)
	wantBuffer, _ := createNewBufferWithData([][]string{{"a", "b", "c"}, {"10", "8", "9"}, {"2", "2", "3"}, {"4", "5", "6"}}, true)
	tests := []struct {
		name string
		b    *Buffer
		args args
		want *Buffer
	}{
		{"test", testBuffer, args{colIndex: 0, rev: false}, wantBuffer},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.sortByStr(tt.args.colIndex, tt.args.rev)
			if !reflect.DeepEqual(tt.b, tt.want) {
				t.Errorf("Buffer_sortByStr() = %v, want %v", tt.b, tt.want)
			}
		})
	}
}

func TestBuffer_sortByNum(t *testing.T) {
	type args struct {
		colIndex int
		rev      bool
	}
	testBuffer, _ := createNewBufferWithData([][]string{{"a", "b", "c"}, {"5", "2", "3"}, {"4", "5", "6"}, {"10", "8", "9"}}, true)
	wantBuffer, _ := createNewBufferWithData([][]string{{"a", "b", "c"}, {"4", "5", "6"}, {"5", "2", "3"}, {"10", "8", "9"}}, true)
	tests := []struct {
		name string
		b    *Buffer
		args args
		want *Buffer
	}{
		{"test", testBuffer, args{colIndex: 0, rev: false}, wantBuffer},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.sortByNum(tt.args.colIndex, tt.args.rev)
			if !reflect.DeepEqual(tt.b, tt.want) {
				t.Errorf("Buffer_sortByNum() = %v, want %v", tt.b, tt.want)
			}
		})
	}
}

func TestBuffer_transpose(t *testing.T) {
	testBuffer, _ := createNewBufferWithData([][]string{{"a", "b", "c"}, {"1", "2", "3"}, {"4", "5", "6"}}, true)
	wantBuffer, _ := createNewBufferWithData([][]string{{"a", "1", "4"}, {"b", "2", "5"}, {"c", "3", "6"}}, true)
	tests := []struct {
		name string
		b    *Buffer
		want *Buffer
	}{
		{"test", testBuffer, wantBuffer},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.transpose()
			if !reflect.DeepEqual(tt.b, tt.want) {
				t.Errorf("Buffer_transpose() = %v, want %v", tt.b, tt.want)
			}
		})
	}
}

func TestBuffer_getCol(t *testing.T) {
	type args struct {
		i int
	}
	testBuffer, _ := createNewBufferWithData([][]string{{"a", "b", "c"}, {"1", "2", "3"}, {"4", "5", "6"}}, true)
	tests := []struct {
		name string
		b    *Buffer
		args args
		want []string
	}{
		{"test", testBuffer, args{i: 0}, []string{"a", "1", "4"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.getCol(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Buffer.getCol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuffer_selectBySearch(t *testing.T) {
	type args struct {
		s string
	}
	testBuffer, _ := createNewBufferWithData([][]string{{"a", "b", "c"}, {"1", "2", "3"}, {"4", "1", "6"}}, true)
	tests := []struct {
		name string
		b    *Buffer
		args args
		want [][]int
	}{
		{"test", testBuffer, args{s: "1"}, [][]int{{1, 0}, {2, 1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b.selectBySearch(tt.args.s)
			if got := b.selectedCell; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkVisible() = %v, want %v", got, tt.want)
			}
		})
	}
}
