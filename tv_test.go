package main

import (
	"os"
	"testing"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"main run"},
	}
	os.Args = []string{"data/test.csv"}
	debug = true
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
