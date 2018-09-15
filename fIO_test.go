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
		{"load file into buffer", args{"data/test.csv", createNewBuffer()}, false, [][]string{[]string{"A", "B", "C"}, []string{"1", "2222", "3"}, []string{"2", "1628", "3"}}},
		{"load file into buffer", args{"data/test.tsv", createNewBuffer()}, false, [][]string{[]string{"A", "B", "C"}, []string{"1", "2222", "3"}, []string{"2", "1628", "3"}}},
		{"load compressed file into buffer", args{"data/test.csv.gz", createNewBuffer()}, false, [][]string{[]string{"A", "B", "C"}, []string{"1", "2222", "3"}, []string{"2", "1628", "3"}}},
		{"load compressed file into buffer", args{"data/test.tsv.gz", createNewBuffer()}, false, [][]string{[]string{"A", "B", "C"}, []string{"1", "2222", "3"}, []string{"2", "1628", "3"}}},
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
