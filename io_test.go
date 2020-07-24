package main

import (
	"reflect"
	"testing"
)

func Test_loadFileToBuffer(t *testing.T) {
	type args struct {
		fn string
		b  *Buffer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"load tsv file", args{fn: "./data/test.tsv", b: createNewBuffer()}, false},
		{"load csv file", args{fn: "./data/test.csv", b: createNewBuffer()}, false},
		{"load gzip file", args{fn: "./data/test.csv.gz", b: createNewBuffer()}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := loadFileToBuffer(tt.args.fn, tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("loadFileToBuffer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_skipLine(t *testing.T) {
	type args struct {
		line string
		sy   []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test", args{"@some line fo filtrate", []string{"@some"}}, true},
		{"test", args{"@some line fo filtrate", []string{"some"}}, false},
		{"test", args{"@some line fo filtrate", []string{"@some", "some"}}, true},
		{"test", args{"some @some line fo filtrate", []string{"@some", "some"}}, true},
		{"test", args{"some @some line fo filtrate", []string{"@some"}}, false},
		{"test", args{"some @some line fo filtrate", []string{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := skipLine(tt.args.line, tt.args.sy); got != tt.want {
				t.Errorf("skipLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getVisCol(t *testing.T) {
	type args struct {
		showNumL []int
		hideNumL []int
		colLen   int
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{"set show argument", args{[]int{1, 2, 5}, []int{}, 6}, []int{0, 1, 4}, false},
		{"set hide argument", args{[]int{}, []int{1, 2, 5}, 6}, []int{2, 3, 5}, false},
		{"do not set any argument", args{[]int{}, []int{}, 6}, []int{0., 1, 2, 3, 4, 5}, false},
		{"argument error", args{[]int{1, 2, 3}, []int{1, 2, 3}, 6}, nil, true},
		{"column does not exist", args{[]int{1, 2, 3, 7}, []int{}, 6}, nil, true},
		{"column does not exist", args{[]int{}, []int{1, 2, 3, 7}, 6}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getVisCol(tt.args.showNumL, tt.args.hideNumL, tt.args.colLen)
			if (err != nil) != tt.wantErr {
				t.Errorf("getVisCol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getVisCol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkVisible(t *testing.T) {
	type args struct {
		showNumL []int
		hideNumL []int
		col      int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"set show argument", args{[]int{1, 3, 5}, []int{}, 5}, false, false},
		{"set show argument", args{[]int{1, 3, 5}, []int{}, 4}, true, false},
		{"set show argument", args{[]int{1, 3, 5}, []int{}, 3}, false, false},
		{"set show argument", args{[]int{1, 3, 5}, []int{}, 2}, true, false},
		{"set show argument", args{[]int{1, 3, 5}, []int{}, 1}, false, false},
		{"set show argument", args{[]int{1, 3, 5}, []int{}, 0}, true, false},
		{"set hide argument", args{[]int{}, []int{1, 3, 5}, 5}, true, false},
		{"set hide argument", args{[]int{}, []int{1, 3, 5}, 4}, false, false},
		{"set hide argument", args{[]int{}, []int{1, 3, 5}, 3}, true, false},
		{"set hide argument", args{[]int{}, []int{1, 3, 5}, 2}, false, false},
		{"set hide argument", args{[]int{}, []int{1, 3, 5}, 1}, true, false},
		{"set hide argument", args{[]int{}, []int{1, 3, 5}, 0}, false, false},
		{"argument error", args{[]int{1, 3, 5}, []int{1, 3, 5}, 4}, false, true},
		{"do not set any argument", args{[]int{}, []int{}, 3}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkVisible(tt.args.showNumL, tt.args.hideNumL, tt.args.col)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkVisible() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkVisible() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lineCSVParse(t *testing.T) {
	type args struct {
		s   string
		sep rune
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"csv line test", args{"a,b,c,d", ','}, []string{"a", "b", "c", "d"}, false},
		{"tsv line test", args{"a	b	c	d", '	'}, []string{"a", "b", "c", "d"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lineCSVParse(tt.args.s, tt.args.sep)
			if (err != nil) != tt.wantErr {
				t.Errorf("lineCSVParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lineCSVParse() = %v, want %v", got, tt.want)
			}
		})
	}
}
