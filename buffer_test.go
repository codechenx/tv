package main

import (
	"reflect"
	"testing"
)

func Test_createNewBuffer(t *testing.T) {
	tests := []struct {
		name string
		want *Buffer
	}{
		{"createNewBuffer", &Buffer{name: "", sep: "", header: 0, vHRN: 0, vHCN: 0, rowNum: 0, colNum: 0, cont: [][]string{}}},
	}
	for _, tt := range tests {
		if got := createNewBuffer(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. createNewBuffer() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestBuffer_contAppend(t *testing.T) {
	type fields struct {
		name   string
		sep    string
		header int
		cont   [][]string
		rowNum int
		colNum int
	}
	type args struct {
		s   string
		sep string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    [][]string
		wantErr bool
	}{
		{"Add string to cont", fields{"", "", 2, [][]string{}, 1, 4}, args{"some,thing,other,thing", ","}, [][]string{[]string{"some", "thing", "other", "thing"}}, false},
		{"Add string to cont", fields{"", "", 2, [][]string{[]string{"some", "thing", "other", "thing"}}, 2, 4}, args{"some,thing,other,thing", ","}, [][]string{[]string{"some", "thing", "other", "thing"}, []string{"some", "thing", "other", "thing"}}, false},
		{"Add string to cont", fields{"", "", 2, [][]string{}, 1, 4}, args{"some	thing	other	thing", "	"}, [][]string{[]string{"some", "thing", "other", "thing"}}, false},
		{"Add string to cont", fields{"", "", 2, [][]string{[]string{"some", "thing", "other", "thing"}}, 2, 4}, args{"some	thing	other	thing", "	"}, [][]string{[]string{"some", "thing", "other", "thing"}, []string{"some", "thing", "other", "thing"}}, false},
		{"Add string to cont", fields{"", "", 2, [][]string{[]string{"some", "thing", "other", "thing"}}, 2, 4}, args{"some,thing,other", ","}, [][]string{[]string{"some", "thing", "other", "thing"}}, true},
	}
	for _, tt := range tests {
		b := &Buffer{
			name:   tt.fields.name,
			sep:    tt.fields.sep,
			header: tt.fields.header,
			cont:   tt.fields.cont,
			rowNum: tt.fields.rowNum,
			colNum: tt.fields.colNum,
		}

		err := b.contAppend(tt.args.s, tt.args.sep)
		if got := b.cont; !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Buffer.contAppend(string,string) = %v, want = %v", tt.name, got, tt.want)
		}
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Buffer.contAppend() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestBuffer_addVirHeader(t *testing.T) {
	type fields struct {
		name   string
		sep    string
		header int
		cont   [][]string
		rowNum int
		colNum int
	}
	tests := []struct {
		name   string
		fields fields
		want   [][]string
	}{
		{"Add virtual header to cont", fields{"", "", 2, [][]string{[]string{"some", "thing", "other", "thing"}, []string{"some", "thing", "other", "thing"}}, 2, 4}, [][]string{[]string{"some", "thing", "other", "thing"}, []string{"some", "thing", "other", "thing"}}},
		{"Add virtual header to cont", fields{"", "", 1, [][]string{[]string{"some", "thing", "other", "thing"}, []string{"some", "thing", "other", "thing"}}, 2, 4}, [][]string{[]string{"1", "2", "3", "4"}, []string{"some", "thing", "other", "thing"}, []string{"some", "thing", "other", "thing"}}},
		{"Add virtual header to cont", fields{"", "", 0, [][]string{[]string{"some", "thing", "other", "thing"}, []string{"some", "thing", "other", "thing"}}, 2, 4}, [][]string{[]string{"#", "some", "thing", "other", "thing"}, []string{"1", "some", "thing", "other", "thing"}}},
		{"Add virtual header to cont", fields{"", "", -1, [][]string{[]string{"some", "thing", "other", "thing"}, []string{"some", "thing", "other", "thing"}}, 2, 4}, [][]string{[]string{"#", "1", "2", "3", "4"}, []string{"1", "some", "thing", "other", "thing"}, []string{"2", "some", "thing", "other", "thing"}}},
	}
	for _, tt := range tests {
		b := Buffer{
			name:   tt.fields.name,
			sep:    tt.fields.sep,
			header: tt.fields.header,
			cont:   tt.fields.cont,
			rowNum: tt.fields.rowNum,
			colNum: tt.fields.colNum,
		}
		b.addVirHeader()
		if got := b.cont; !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Buffer.addVirHeader() = %v, want = %v", tt.name, got, tt.want)
		}
	}
}
