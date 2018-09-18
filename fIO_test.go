package main

import (
	"reflect"
	"testing"
)

func Test_loadFile(t *testing.T) {
	type args struct {
		fn string
		b  *Buffer
	}

	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantCont [][]string
	}{
		{"load file into buffer", args{"data/test.csv", createNewBuffer()}, false, [][]string{{"A", "B", "C"}, {"1", "2222", "3"}, {"2", "1628", "3"}}},
		{"load file into buffer", args{"data/test.tsv", createNewBuffer()}, false, [][]string{{"A", "B", "C"}, {"1", "2222", "3"}, {"2", "1628", "3"}}},
		{"load compressed file into buffer", args{"data/test.csv.gz", createNewBuffer()}, false, [][]string{{"A", "B", "C"}, {"1", "2222", "3"}, {"2", "1628", "3"}}},
		{"load compressed file into buffer", args{"data/test.tsv.gz", createNewBuffer()}, false, [][]string{{"A", "B", "C"}, {"1", "2222", "3"}, {"2", "1628", "3"}}},
		{"file does not exist", args{"file_does_not_exits", createNewBuffer()}, true, [][]string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := loadFile(tt.args.fn, tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("loadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		if got := tt.args.b.cont; !reflect.DeepEqual(got, tt.wantCont) {
			t.Errorf("%q. loadFile(path,buffer) = %v, wantCont = %v", tt.name, got, tt.wantCont)
		}
	}
}

func Test_exists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"file exists?", args{"path"}, false},
		{"file exists?", args{"tv.go"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exists(tt.args.path); got != tt.want {
				t.Errorf("exists() = %v, want %v", got, tt.want)
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
		{"some line fo test", args{"@some line fo filtrate", []string{"@some"}}, true},
		{"some line fo test", args{"@some line fo filtrate", []string{"some"}}, false},
		{"some line fo test", args{"@some line fo filtrate", []string{"@some", "some"}}, true},
		{"some line fo test", args{"some @some line fo filtrate", []string{"@some", "some"}}, true},
		{"some line fo test", args{"some @some line fo filtrate", []string{"@some"}}, false},
		{"some line fo test", args{"some @some line fo filtrate", []string{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := skipLine(tt.args.line, tt.args.sy); got != tt.want {
				t.Errorf("skipLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compressed(t *testing.T) {
	type args struct {
		fn string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"compressed test", args{"file.gz"}, true},
		{"compressed test", args{"file.gz.gz"}, true},
		{"compressed test", args{"file.gz.txt"}, false},
		{"compressed test", args{"file.g"}, false},
		{"compressed test", args{"gz."}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compressed(tt.args.fn); got != tt.want {
				t.Errorf("compressed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deterSep(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"deterSep test", args{"some,thing,other,thing"}, ",", false},
		{"deterSep test", args{"some	thing@other@thing"}, "	", false},
		{"deterSep test", args{"some	thing	other	thing"}, "\t", false},
		{"deterSep test", args{"some@thing@other@thing"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := deterSep(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("deterSep() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("deterSep() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fScanner(t *testing.T) {
	type args struct {
		fn   string
		comp bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"fScanner scan not", args{"data/test.tsv", false}, false},
		{"fScanner scan not", args{"data/test.tsv.gz", true}, false},
		{"fScanner scan not", args{"data/test.csv", false}, false},
		{"fScanner scan not", args{"data/test.csv.gz", true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fScanner(tt.args.fn, tt.args.comp)
			if (err != nil) != tt.wantErr {
				t.Errorf("fScanner() error = %v, wantErr %v", err, tt.wantErr)
				return
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
