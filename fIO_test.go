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
		{"load file", args{"data/test.csv", createNewBuffer()}, false, [][]string{[]string{"A", "B", "C"}, []string{"1", "2222", "3"}, []string{"2", "1628", "3"}}},
		{"load file", args{"data/test.tsv", createNewBuffer()}, false, [][]string{[]string{"A", "B", "C"}, []string{"1", "2222", "3"}, []string{"2", "1628", "3"}}},
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
